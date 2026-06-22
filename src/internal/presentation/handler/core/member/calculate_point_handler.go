package member

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"go-oapi-aidd/internal/core/member/usecase"
	"go-oapi-aidd/internal/presentation/gen"
)

type CalculatePointHandler struct {
	usecase usecase.CalculatePointUsecase
}

func NewCalculatePointHandler(usecase usecase.CalculatePointUsecase) *CalculatePointHandler {
	return &CalculatePointHandler{
		usecase: usecase,
	}
}

func (h *CalculatePointHandler) CalculateMemberPoint(
	ctx context.Context,
	request gen.CalculateMemberPointRequestObject,
) (gen.CalculateMemberPointResponseObject, error) {
	if request.Body == nil {
		return badRequestResponse(), nil
	}

	input := usecase.CalculatePointInput{
		MemberID:       request.MemberId.String(),
		PurchaseAmount: request.Body.PurchaseAmount,
	}

	output, err := h.usecase.Execute(ctx, input)
	if err != nil {
		if errors.Is(err, usecase.ErrMemberNotFound) {
			return gen.CalculateMemberPoint404JSONResponse{
				Code:    "MEMBER_NOT_FOUND",
				Message: "member not found",
			}, nil
		}
		return gen.CalculateMemberPoint500JSONResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		}, nil
	}

	memberID, err := uuid.Parse(output.MemberID)
	if err != nil {
		return gen.CalculateMemberPoint500JSONResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		}, nil
	}

	return gen.CalculateMemberPoint200JSONResponse{
		MemberId:       memberID,
		PurchaseAmount: output.PurchaseAmount,
		GrantedPoint:   output.GrantedPoint,
	}, nil
}

func badRequestResponse() gen.CalculateMemberPoint400JSONResponse {
	return gen.CalculateMemberPoint400JSONResponse{
		Code:    "BAD_REQUEST",
		Message: "bad request",
	}
}
