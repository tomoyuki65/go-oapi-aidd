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
	"github.com/uptrace/bun"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/infrastructure/database/schema"
	"go-oapi-aidd/internal/infrastructure/logger"
	"go-oapi-aidd/internal/presentation/router"
)

const tanakaMemberID = "11111111-1111-1111-1111-111111111111"

func TestMemberPointCalculation(t *testing.T) {
	ctx := context.Background()
	db, err := database.NewBunDB()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = cleanupMembers(ctx, db, tanakaMemberID)
		db.Close()
	})
	require.NoError(t, cleanupMembers(ctx, db, tanakaMemberID))
	require.NoError(t, insertMembers(ctx, db, schema.Member{
		ID:   tanakaMemberID,
		Name: "田中",
		Rank: "bronze",
	}))

	container := di.NewContainer(db, logger.NewSlogLogger())
	r := router.NewRouter(container)

	t.Run("テスト用会員IDへのHTTPリクエストでポイント計算結果を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/"+tanakaMemberID+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"memberId":"`+tanakaMemberID+`","purchaseAmount":5000,"grantedPoint":50}`, res.Body.String())
	})

	t.Run("存在しない会員IDへのHTTPリクエストで404を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/99999999-9999-9999-9999-999999999999/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, `{"code":"MEMBER_NOT_FOUND","message":"member not found"}`, res.Body.String())
	})

	t.Run("UUID形式ではないmemberIdへのHTTPリクエストで400を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/not-a-uuid/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{"message":"Bad Request"}`, res.Body.String())
	})

	t.Run("purchaseAmountが負数の場合に400を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/"+tanakaMemberID+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":-1}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{"message":"Bad Request"}`, res.Body.String())
	})

	t.Run("purchaseAmountが1000000000の場合に400を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/"+tanakaMemberID+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":1000000000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{"message":"Bad Request"}`, res.Body.String())
	})

	t.Run("purchaseAmountが未指定の場合に400を返すこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/members/"+tanakaMemberID+"/point-calculations", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{"message":"Bad Request"}`, res.Body.String())
	})
}

func insertMembers(ctx context.Context, db *bun.DB, members ...schema.Member) error {
	_, err := db.NewInsert().Model(&members).On("CONFLICT (id) DO NOTHING").Exec(ctx)
	return err
}

func cleanupMembers(ctx context.Context, db *bun.DB, ids ...string) error {
	_, err := db.NewDelete().Model((*schema.Member)(nil)).Where("id IN (?)", bun.List(ids)).Exec(ctx)
	return err
}
