//go:build unit

package member

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	openapiTypes "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-oapi-aidd/internal/core/member/usecase"
	mockUsecase "go-oapi-aidd/internal/core/member/usecase/mock_usecase"
	"go-oapi-aidd/internal/presentation/gen"
)

type recordingServer struct {
	called bool
}

func (s *recordingServer) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *recordingServer) CalculateMemberPoint(w http.ResponseWriter, r *http.Request, memberId openapiTypes.UUID) {
	s.called = true
	w.WriteHeader(http.StatusNoContent)
}

func TestCalculatePointHandler_CalculateMemberPoint(t *testing.T) {
	memberID := openapiTypes.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}

	t.Run("path parameterのmemberIdをusecase入力へ変換すること", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{
				MemberID:       memberID.String(),
				PurchaseAmount: 5000,
			}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 5000, GrantedPoint: 50}, nil)

		handler := NewCalculatePointHandler(mockUsecase)
		_, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 5000},
		})

		assert.NoError(t, err)
	})

	t.Run("request bodyのpurchaseAmountをusecase入力へ変換すること", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{
				MemberID:       memberID.String(),
				PurchaseAmount: 999999999,
			}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 999999999, GrantedPoint: 19999998}, nil)

		handler := NewCalculatePointHandler(mockUsecase)
		_, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 999999999},
		})

		assert.NoError(t, err)
	})

	t.Run("purchaseAmountの境界値0円をusecase入力へ変換すること", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{MemberID: memberID.String(), PurchaseAmount: 0}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 0, GrantedPoint: 0}, nil)

		handler := NewCalculatePointHandler(mockUsecase)
		_, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 0},
		})

		assert.NoError(t, err)
	})

	t.Run("purchaseAmountの境界値999999999円をusecase入力へ変換すること", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{MemberID: memberID.String(), PurchaseAmount: 999999999}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 999999999, GrantedPoint: 19999998}, nil)

		handler := NewCalculatePointHandler(mockUsecase)
		_, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 999999999},
		})

		assert.NoError(t, err)
	})

	t.Run("purchaseAmountが負数の場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)

		handler := NewCalculatePointHandler(mockUsecase)
		res, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: -1},
		})

		jsonRes, ok := res.(gen.CalculateMemberPoint400JSONResponse)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "BAD_REQUEST", jsonRes.Code)
	})

	t.Run("purchaseAmountが1000000000の場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)

		handler := NewCalculatePointHandler(mockUsecase)
		res, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 1000000000},
		})

		jsonRes, ok := res.(gen.CalculateMemberPoint400JSONResponse)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "BAD_REQUEST", jsonRes.Code)
	})

	t.Run("purchaseAmountが未指定の場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)

		handler := NewCalculatePointHandler(mockUsecase)
		res, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     nil,
		})

		jsonRes, ok := res.(gen.CalculateMemberPoint400JSONResponse)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "BAD_REQUEST", jsonRes.Code)
	})

	t.Run("200レスポンスのJSON shapeを検証すること", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 5000, GrantedPoint: 50}, nil)

		handler := NewCalculatePointHandler(mockUsecase)
		res, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 5000},
		})

		jsonRes, ok := res.(gen.CalculateMemberPoint200JSONResponse)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, memberID, jsonRes.MemberId)
		assert.Equal(t, 5000, jsonRes.PurchaseAmount)
		assert.Equal(t, 50, jsonRes.GrantedPoint)
	})

	t.Run("会員未存在時に404とMEMBER_NOT_FOUNDを返すこと", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{}, usecase.ErrMemberNotFound)

		handler := NewCalculatePointHandler(mockUsecase)
		res, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 5000},
		})

		jsonRes, ok := res.(gen.CalculateMemberPoint404JSONResponse)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "MEMBER_NOT_FOUND", jsonRes.Code)
	})

	t.Run("内部エラー時に500とINTERNAL_SERVER_ERRORを返すこと", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{}, usecase.ErrInternal)

		handler := NewCalculatePointHandler(mockUsecase)
		res, err := handler.CalculateMemberPoint(context.Background(), gen.CalculateMemberPointRequestObject{
			MemberId: memberID,
			Body:     &gen.CalculateMemberPointJSONRequestBody{PurchaseAmount: 5000},
		})

		jsonRes, ok := res.(gen.CalculateMemberPoint500JSONResponse)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", jsonRes.Code)
	})
}

func TestCalculatePointHandler_InvalidMemberID(t *testing.T) {
	t.Run("UUID形式ではないmemberIdの場合に400を返しusecase相当の処理を呼び出さないこと", func(t *testing.T) {
		server := &recordingServer{}
		handler := gen.Handler(server)
		req := httptest.NewRequest(http.MethodPost, "/members/not-a-uuid/point-calculations", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.False(t, server.called)
	})
}
