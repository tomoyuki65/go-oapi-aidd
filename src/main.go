package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/infrastructure/logger"
	"go-oapi-aidd/internal/infrastructure/observability"
	"go-oapi-aidd/internal/presentation/router"
)

func main() {
	// コンテキスト設定
	ctx := context.Background()

	// OpenTelemetryのトレース出力先設定
	otelExporterType := os.Getenv("OTEL_EXPORTER_TYPE")
	shutdownTracer := observability.InitTracer(ctx, otelExporterType)
	defer func() {
		_ = shutdownTracer(context.Background())
	}()

	// DB取得
	db, err := database.NewBunDB()
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// ロガー取得
	logger := logger.NewSlogLogger()

	// DIコンテナ取得
	container := di.NewContainer(db, logger)

	// ルーティング設定の取得
	r := router.NewRouter(container)

	otelHandler := otelhttp.NewHandler(
		r,
		"go-oapi-aidd",
		otelhttp.WithSpanNameFormatter(func(operation string, req *http.Request) string {
			return req.Method + " " + req.URL.Path
		}),
	)

	// ポート番号の設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	startPort := fmt.Sprintf(":%s", port)

	// サーバー設定
	srv := &http.Server{
		Addr:              startPort,
		Handler:           otelHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// サーバー起動
	logger.Info(false, ctx, "start server go-oapi-aidd")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(false, ctx, err.Error())
		os.Exit(1)
	}
}
