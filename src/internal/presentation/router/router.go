package router

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	nethttpMiddleware "github.com/oapi-codegen/nethttp-middleware"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/presentation/gen"
	"go-oapi-aidd/internal/presentation/handler"
)

// リクエストのタイムアウト設定値取得関数
func getRequestTimeout() time.Duration {
	// 環境変数「REQUEST_TIMEOUT_SECONDS」の値を取得
	requestTimeoutSecondsStr := os.Getenv("REQUEST_TIMEOUT_SECONDS")
	if requestTimeoutSecondsStr == "" {
		requestTimeoutSecondsStr = "10"
	}

	// INT型に変換
	requestTimeoutSecondsInt, err := strconv.Atoi(requestTimeoutSecondsStr)
	if err != nil {
		panic(err)
	}

	// Duration型の秒数で返す
	return time.Duration(requestTimeoutSecondsInt) * time.Second
}

// 許可されたオリジン設定値取得関数
func getAllowedOrigins() []string {
	// 環境変数「CORS_ALLOWED_ORIGIN」の値を取得
	allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return []string{
		allowedOrigin,
	}
}

// OpenAPIを利用したバリデーションチェック追加用関数
func addOapiValidation(path string, r chi.Router) {
	// Swagger（OpenAPI定義）取得
	swagger, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}

	// Servers定義を対象のpathで上書き
	swagger.Servers = openapi3.Servers{
		&openapi3.Server{
			URL: path,
		},
	}

	// バリデーション設定
	validator := nethttpMiddleware.OapiRequestValidatorWithOptions(
		swagger,
		&nethttpMiddleware.Options{
			SilenceServersWarning: true,
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_ = json.NewEncoder(w).Encode(map[string]any{
					// "code":    statusCode,
					"message": message,
				})
			},
		},
	)

	// ミドルウェアにvalidatorを追加
	r.Use(validator)
}

func NewRouter(container *di.Container) *chi.Mux {
	r := chi.NewRouter()

	// ミドルウェアの設定
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(getRequestTimeout()))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: getAllowedOrigins(),
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	// テスト実行時以外にロガーを追加
	if env := os.Getenv("ENV"); env != "testing" {
		r.Use(middleware.Logger)
	}

	// 「/api/v1」のルーティング設定
	handlerV1 := handler.NewHandlerV1(container)
	strictHandlerV1 := gen.NewStrictHandler(handlerV1, nil)
	apiV1 := "/api/v1"
	r.Route(apiV1, func(r chi.Router) {
		addOapiValidation(apiV1, r)
		gen.HandlerFromMux(strictHandlerV1, r)
	})

	return r
}
