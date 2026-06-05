//go:build unit

package healthcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/presentation/gen"
	mockLogger "go-oapi-aidd/internal/shared/logger/mock_logger"
	mockService "go-oapi-aidd/internal/supporting/healthcheck/mock_service"
)

func TestHealthcheckHandler_Healthcheck(t *testing.T) {
	// DB取得
	db, err := database.NewBunDB()
	if err != nil {
		t.Fatal("failed to connect database")
	}

	// ロガーのモック化
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mockLogger.NewMockLogger(ctrl)

	// ユースケースのモック化
	mockService := mockService.NewMockService(ctrl)

	t.Run("正常終了すること", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().Execute(gomock.Any()).Return(nil)

		// DIコンテナ設定
		container := di.NewContainer(db, mockLogger)

		// ハンドラー設定
		handler := NewHealthcheckHandler(container, mockService)

		// テスト実行
		res, err := handler.Healthcheck(
			context.Background(),
			gen.HealthcheckRequestObject{},
		)

		// レスポンス結果の変換
		jsonRes, ok := res.(gen.Healthcheck200JSONResponse)

		// 検証
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "OK", jsonRes.Message)
	})

	t.Run("異常終了すること", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().Execute(gomock.Any()).Return(assert.AnError)

		// DIコンテナ設定
		container := di.NewContainer(db, mockLogger)

		// ハンドラー設定
		handler := NewHealthcheckHandler(container, mockService)

		// テスト実行
		res, err := handler.Healthcheck(
			context.Background(),
			gen.HealthcheckRequestObject{},
		)

		// レスポンス結果の変換
		jsonRes, ok := res.(gen.Healthcheck500JSONResponse)

		// 検証
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "Internal Server Error", jsonRes.Message)
	})
}
