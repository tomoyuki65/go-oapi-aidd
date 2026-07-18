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

const (
	httpShutdownTimeout   = 10 * time.Second
	tracerShutdownTimeout = 5 * time.Second
)

func main() {
	// 全体の共通コンテキスト設定
	rootCtx := context.Background()

	// SIGINT / SIGTERMを受け取るとキャンセルされるContext
	signalCtx, stopSignal := signal.NotifyContext(
		rootCtx,
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignal()

	// OpenTelemetryのトレース出力先設定
	otelExporterType := os.Getenv("OTEL_EXPORTER_TYPE")
	shutdownTracer := observability.InitTracer(rootCtx, otelExporterType)

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
		logger.Info(false, rootCtx, "start server go-oapi-aidd")

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	// 終了シグナルまたはサーバーエラーの待機
	select {
	case <-signalCtx.Done():
		logger.Info(false, rootCtx, "shutdown signal received")

		// シグナル監視を解除（これ以降に再度Ctrl+Cを押すと、OSのデフォルト動作で終了）
		stopSignal()

	case err := <-serverErrCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(false, rootCtx, err.Error())
		}
	}

	// // Shutdown用のコンテキスト設定
	// shutdownCtx, cancel := context.WithTimeout(
	// 	context.Background(),
	// 	10*time.Second,
	// )
	// defer cancel()

	httpShutdownCtx, cancelHTTPShutdown := context.WithTimeout(
		rootCtx,
		httpShutdownTimeout,
	)

	logger.Info(false, rootCtx, "shutting down HTTP server")

	// 新規リクエストの受付を停止し、処理中のリクエストの完了待ち
	if err := srv.Shutdown(httpShutdownCtx); err != nil {
		logger.Error(
			false,
			rootCtx,
			fmt.Sprintf("failed to gracefully shutdown HTTP server: %v", err),
		)

		// タイムアウトなどでShutdownできなかった場合は強制終了
		if closeErr := srv.Close(); closeErr != nil {
			logger.Error(
				false,
				rootCtx,
				fmt.Sprintf("failed to close HTTP server: %v", closeErr),
			)
		}
	}

	cancelHTTPShutdown()

	// HTTPリクエスト終了後、依存リソースを閉じる
	if err := db.Close(); err != nil {
		logger.Error(
			false,
			rootCtx,
			fmt.Sprintf("failed to close database: %v", err),
		)
	}

	// HTTPサーバーの停止に時間を使い切っていてもTracerを停止できるよう、
	// Tracer用には別のContextを作成する。
	tracerShutdownCtx, cancelTracerShutdown := context.WithTimeout(
		rootCtx,
		tracerShutdownTimeout,
	)
	defer cancelTracerShutdown()

	if err := shutdownTracer(tracerShutdownCtx); err != nil {
		logger.Error(
			false,
			rootCtx,
			fmt.Sprintf("failed to shutdown tracer: %v", err),
		)
	}

	logger.Info(false, rootCtx, "server shutdown completed")
}
