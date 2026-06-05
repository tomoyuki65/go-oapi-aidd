//go:build unit

package healthcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-oapi-aidd/internal/infrastructure/database"
	mockLogger "go-oapi-aidd/internal/shared/logger/mock_logger"
)

func TestService_Execute(t *testing.T) {
	// DB取得
	db, err := database.NewBunDB()
	if err != nil {
		t.Fatal("failed to connect database")
	}

	// ロガーのモック化
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mockLogger.NewMockLogger(ctrl)

	t.Run("正常終了すること", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		// サービス定義
		service := NewService(db, mockLogger)

		// テスト実行
		err := service.Execute(context.Background())

		// 検証
		assert.NoError(t, err)
	})

	t.Run("異常終了すること", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		// サービス定義
		db.Close()
		service := NewService(db, mockLogger)

		// テスト実行
		err := service.Execute(context.Background())

		// 検証
		assert.Error(t, err)
	})
}
