//go:build unit

package member

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	openapiTypes "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-oapi-aidd/internal/core/member/usecase"
	mockUsecase "go-oapi-aidd/internal/core/member/usecase/mock_usecase"
	"go-oapi-aidd/internal/presentation/gen"
)

type strictServer struct {
	handler *CalculatePointHandler
}

func (s strictServer) Healthcheck(
	ctx context.Context,
	request gen.HealthcheckRequestObject,
) (gen.HealthcheckResponseObject, error) {
	return gen.Healthcheck500JSONResponse{Message: "not implemented"}, nil
}

func (s strictServer) CalculateMemberPoint(
	ctx context.Context,
	request gen.CalculateMemberPointRequestObject,
) (gen.CalculateMemberPointResponseObject, error) {
	return s.handler.CalculateMemberPoint(ctx, request)
}

func newTestHandler(usecase usecase.CalculatePointUsecase) http.Handler {
	handler := NewCalculatePointHandler(usecase)
	r := chi.NewRouter()
	gen.HandlerFromMux(gen.NewStrictHandler(strictServer{handler: handler}, nil), r)
	return r
}

func TestCalculatePointHandler_CalculateMemberPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := mockUsecase.NewMockCalculatePointUsecase(ctrl)
	handler := newTestHandler(mockUsecase)
	memberID := openapiTypes.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}

	t.Run("path parameterのmemberIdをusecase入力へ変換すること", func(t *testing.T) {
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{
				MemberID:       memberID.String(),
				PurchaseAmount: 5000,
			}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 5000, GrantedPoint: 50}, nil)

		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("request bodyのpurchaseAmountをusecase入力へ変換すること", func(t *testing.T) {
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{
				MemberID:       memberID.String(),
				PurchaseAmount: 999999999,
			}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 999999999, GrantedPoint: 19999998}, nil)

		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":999999999}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("purchaseAmountの境界値0円をusecase入力へ変換すること", func(t *testing.T) {
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{MemberID: memberID.String(), PurchaseAmount: 0}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 0, GrantedPoint: 0}, nil)

		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":0}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("purchaseAmountの境界値999999999円をusecase入力へ変換すること", func(t *testing.T) {
		mockUsecase.EXPECT().
			Execute(gomock.Any(), usecase.CalculatePointInput{MemberID: memberID.String(), PurchaseAmount: 999999999}).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 999999999, GrantedPoint: 19999998}, nil)

		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":999999999}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("UUID形式ではないmemberIdの場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/members/not-a-uuid/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("purchaseAmountが負数の場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":-1}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("purchaseAmountが1000000000の場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":1000000000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("purchaseAmountが未指定の場合に400を返しusecase mockが呼び出されないこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("200レスポンスのJSON shapeを検証すること", func(t *testing.T) {
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{MemberID: memberID.String(), PurchaseAmount: 5000, GrantedPoint: 50}, nil)

		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"memberId":"`+memberID.String()+`","purchaseAmount":5000,"grantedPoint":50}`, res.Body.String())
	})

	t.Run("会員未存在時に404とMEMBER_NOT_FOUNDを返すこと", func(t *testing.T) {
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{}, usecase.ErrMemberNotFound)

		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, `{"code":"MEMBER_NOT_FOUND","message":"member not found"}`, res.Body.String())
	})

	t.Run("内部エラー時に500とINTERNAL_SERVER_ERRORを返すこと", func(t *testing.T) {
		mockUsecase.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(usecase.CalculatePointOutput{}, usecase.ErrInternal)

		req := httptest.NewRequest(http.MethodPost, "/members/"+memberID.String()+"/point-calculations", bytes.NewBufferString(`{"purchaseAmount":5000}`))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, `{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}`, res.Body.String())
	})
}
