//go:build unit

package member

import (
	"context"
	"testing"

	openapiTypes "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-oapi-aidd/internal/core/member/usecase"
	mockUsecase "go-oapi-aidd/internal/core/member/usecase/mock_usecase"
	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/presentation/gen"
	mockLogger "go-oapi-aidd/internal/shared/logger/mock_logger"
)

func TestCalculatePointHandler_CalculateMemberPoint(t *testing.T) {
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
	mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)

	memberID := openapiTypes.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}

	t.Run("正常終了すること", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{
				MemberID:       memberID.String(),
				PurchaseAmount: 5000,
			}).
			Return(usecase.CalculatePointOutput{
				MemberID:       memberID.String(),
				PurchaseAmount: 5000,
				GrantedPoint:   50,
			}, nil)

		// DIコンテナ設定
		container := di.NewContainer(db, mockLogger)

		// ハンドラー設定
		handler := NewCalculatePointHandler(container, mockUsecase)

		// テスト実行
		res, err := handler.CalculateMemberPoint(
			context.Background(),
			gen.CalculateMemberPointRequestObject{
				MemberId: memberID,
				Body: &gen.CalculateMemberPointJSONRequestBody{
					PurchaseAmount: 5000,
				},
			},
		)

		// レスポンス結果の変換
		jsonRes, ok := res.(gen.CalculateMemberPoint200JSONResponse)

		// 検証
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, memberID, jsonRes.MemberId)
		assert.Equal(t, 5000, jsonRes.PurchaseAmount)
		assert.Equal(t, 50, jsonRes.GrantedPoint)
	})

	t.Run("request bodyのpurchaseAmountをusecase入力へ変換すること", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{
				MemberID:       memberID.String(),
				PurchaseAmount: 999999999,
			}).
			Return(usecase.CalculatePointOutput{
				MemberID:       memberID.String(),
				PurchaseAmount: 999999999,
				GrantedPoint:   19999998,
			}, nil)

		// DIコンテナ設定
		container := di.NewContainer(db, mockLogger)

		// ハンドラー設定
		handler := NewCalculatePointHandler(container, mockUsecase)

		// テスト実行
		res, err := handler.CalculateMemberPoint(
			context.Background(),
			gen.CalculateMemberPointRequestObject{
				MemberId: memberID,
				Body: &gen.CalculateMemberPointJSONRequestBody{
					PurchaseAmount: 999999999,
				},
			},
		)

		// レスポンス結果の変換
		jsonRes, ok := res.(gen.CalculateMemberPoint200JSONResponse)

		// 検証
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, 999999999, jsonRes.PurchaseAmount)
		assert.Equal(t, 19999998, jsonRes.GrantedPoint)
	})

	t.Run("request bodyがnilの場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		// DIコンテナ設定
		container := di.NewContainer(db, mockLogger)

		// ハンドラー設定
		handler := NewCalculatePointHandler(container, mockUsecase)

		// テスト実行
		res, err := handler.CalculateMemberPoint(
			context.Background(),
			gen.CalculateMemberPointRequestObject{
				MemberId: memberID,
				Body:     nil,
			},
		)

		// レスポンス結果の変換
		jsonRes, ok := res.(gen.CalculateMemberPoint400JSONResponse)

		// 検証
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "Bad Request", jsonRes.Message)
	})

	t.Run("会員未存在時に404とMEMBER_NOT_FOUNDを返すこと", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{}, usecase.ErrMemberNotFound)

		// DIコンテナ設定
		container := di.NewContainer(db, mockLogger)

		// ハンドラー設定
		handler := NewCalculatePointHandler(container, mockUsecase)

		// テスト実行
		res, err := handler.CalculateMemberPoint(
			context.Background(),
			gen.CalculateMemberPointRequestObject{
				MemberId: memberID,
				Body: &gen.CalculateMemberPointJSONRequestBody{
					PurchaseAmount: 5000,
				},
			},
		)

		// レスポンス結果の変換
		jsonRes, ok := res.(gen.CalculateMemberPoint404JSONResponse)

		// 検証
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "MEMBER_NOT_FOUND", jsonRes.Code)
		assert.Equal(t, "member not found", jsonRes.Message)
	})

	t.Run("内部エラー時に500と共通エラーレスポンスを返すこと", func(t *testing.T) {
		// モック設定
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{}, usecase.ErrInternal)

		// DIコンテナ設定
		container := di.NewContainer(db, mockLogger)

		// ハンドラー設定
		handler := NewCalculatePointHandler(container, mockUsecase)

		// テスト実行
		res, err := handler.CalculateMemberPoint(
			context.Background(),
			gen.CalculateMemberPointRequestObject{
				MemberId: memberID,
				Body: &gen.CalculateMemberPointJSONRequestBody{
					PurchaseAmount: 5000,
				},
			},
		)

		// レスポンス結果の変換
		jsonRes, ok := res.(gen.CalculateMemberPoint500JSONResponse)

		// 検証
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "Internal Server Error", jsonRes.Message)
	})
}
