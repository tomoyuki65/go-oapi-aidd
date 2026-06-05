package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5/middleware"

	"go-oapi-aidd/internal/shared/logger"
)

// slogの設定
type SlogHandler struct {
	slog.Handler
}

func (h *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	// rをコピー
	newRecord := r.Clone()

	// runtimeからPCを取得して上書き
	pc, _, _, ok := runtime.Caller(4)
	if ok {
		newRecord.PC = pc
	}

	// ミドルウェアで設定したリクエストIDを取得
	requestID := middleware.GetReqID(ctx)

	if requestID != "" {
		newRecord.AddAttrs(
			slog.String("request_id", requestID),
		)
	}

	return h.Handler.Handle(ctx, newRecord)
}

var slogHandler = &SlogHandler{
	slog.NewJSONHandler(os.Stdout, nil),
}

var slogHandlerAddSource = &SlogHandler{
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}),
}

var (
	newLogger          = slog.New(slogHandler)
	newLoggerAddSource = slog.New(slogHandlerAddSource)
)

// ロガーの設定
type slogLogger struct{}

func NewSlogLogger() logger.Logger {
	return &slogLogger{}
}

func (l *slogLogger) Info(addSource bool, ctx context.Context, message string) {
	if env := os.Getenv("ENV"); env != "testing" {
		if addSource {
			newLoggerAddSource.InfoContext(ctx, message)
		} else {
			newLogger.InfoContext(ctx, message)
		}
	}
}

func (l *slogLogger) Warn(addSource bool, ctx context.Context, message string) {
	if env := os.Getenv("ENV"); env != "testing" {
		if addSource {
			newLoggerAddSource.WarnContext(ctx, message)
		} else {
			newLogger.WarnContext(ctx, message)
		}
	}
}

func (l *slogLogger) Error(addSource bool, ctx context.Context, message string) {
	if env := os.Getenv("ENV"); env != "testing" {
		if addSource {
			newLoggerAddSource.ErrorContext(ctx, message)
		} else {
			newLogger.ErrorContext(ctx, message)
		}
	}
}
