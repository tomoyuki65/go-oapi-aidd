//go:build integration

package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/infrastructure/logger"
	"go-oapi-aidd/internal/presentation/router"
)

func TestHealthcheck_OK(t *testing.T) {
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

	// リクエスト作成
	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthcheck", nil)
	res := httptest.NewRecorder()

	// テスト実行
	r.ServeHTTP(res, req)

	// 検証
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, `{"message":"OK"}`, res.Body.String())
	// ciテスト
}
