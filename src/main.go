package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// DB取得
	db, err := database.NewBunDB()
	if err != nil {
		panic("failed to connect database")
	}

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

	// SIGINT（Ctrl+C）またはSIGTERMを受け取るとキャンセルされるコンテキスト設定
	ctx, stop := signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// サーバー設定
	srv := &http.Server{
		Addr:              startPort,
		Handler:           otelHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// サーバーエラー用のチャネル作成
	serverErrCh := make(chan error, 1)

	// ゴルーチンでサーバー起動
	go func() {
		logger.Info(false, ctx, "start server go-oapi-aidd")

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	// 終了シグナルまたはサーバーエラーの待機
	select {
	case <-ctx.Done():
		logger.Info(false, context.Background(), "shutdown signal received")

	case err := <-serverErrCh:
		if err != nil {
			logger.Error(false, context.Background(), err.Error())
		}
	}

	// Shutdown用のコンテキスト設定
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	logger.Info(false, shutdownCtx, "shutting down HTTP server")

	// 新規リクエストの受付を停止し、処理中のリクエストの完了待ち
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			false,
			context.Background(),
			fmt.Sprintf("failed to gracefully shutdown HTTP server: %v", err),
		)

		// タイムアウトなどでShutdownできなかった場合は強制終了
		if closeErr := srv.Close(); closeErr != nil {
			logger.Error(
				false,
				context.Background(),
				fmt.Sprintf("failed to close HTTP server: %v", closeErr),
			)
		}
	}

	// HTTPリクエスト終了後、依存リソースを閉じる
	if err := db.Close(); err != nil {
		logger.Error(
			false,
			context.Background(),
			fmt.Sprintf("failed to close database: %v", err),
		)
	}

	if err := shutdownTracer(shutdownCtx); err != nil {
		logger.Error(
			false,
			context.Background(),
			fmt.Sprintf("failed to shutdown tracer: %v", err),
		)
	}

	logger.Info(false, context.Background(), "server shutdown completed")
}
