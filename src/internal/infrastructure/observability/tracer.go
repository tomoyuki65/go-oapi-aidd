package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

func initTracer(ctx context.Context, exporterType string) func(context.Context) error {
	switch exporterType {
	case "local":
		// 標準出力向けの Trace Exporter を作成
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			panic(err)
		}

		// TracerProvider を構築し、サービス情報と Exporter を設定
		tp := sdkTrace.NewTracerProvider(
			sdkTrace.WithBatcher(exporter),
			sdkTrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName("go-oapi-aidd"),
			)),
		)

		// アプリケーション全体で利用する TracerProvider として登録
		otel.SetTracerProvider(tp)

		return tp.Shutdown

	default:
		// 何もしない関数を返す
		return func(context.Context) error {
			return nil
		}
	}
}
