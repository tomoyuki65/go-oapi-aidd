//go:build integration

package member

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/infrastructure/database/seed"
	"go-oapi-aidd/internal/infrastructure/database/seed/local"
	"go-oapi-aidd/internal/infrastructure/logger"
	"go-oapi-aidd/internal/presentation/router"
)

func TestMemberPointCalculation(t *testing.T) {
	ctx := context.Background()
	db, err := database.NewBunDB()
	require.NoError(t, err)

	memberSeeder := local.NewMemberSeeder()
	t.Cleanup(func() {
		_ = seed.Cleanup(ctx, db, memberSeeder)
		db.Close()
	})
	require.NoError(t, seed.Run(ctx, db, memberSeeder))

	container := di.NewContainer(db, logger.NewSlogLogger())
	r := router.NewRouter(container)

	t.Run("seed済み会員IDへのHTTPリクエストでポイント計算結果を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/"+local.TanakaMemberID+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"memberId":"`+local.TanakaMemberID+`","purchaseAmount":5000,"grantedPoint":50}`, res.Body.String())
	})

	t.Run("存在しない会員IDへのHTTPリクエストで404を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/99999999-9999-9999-9999-999999999999/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, `{"code":"MEMBER_NOT_FOUND","message":"member not found"}`, res.Body.String())
	})

	t.Run("不正なUUIDへのHTTPリクエストで400を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/not-a-uuid/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("範囲外の購入金額へのHTTPリクエストで400を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/"+local.TanakaMemberID+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":1000000000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}
