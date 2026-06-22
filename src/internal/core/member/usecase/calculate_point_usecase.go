package usecase

import (
	"context"
	"errors"
	"fmt"

	"go-oapi-aidd/internal/core/member/domain"
)

var (
	ErrMemberNotFound = errors.New("member not found")
	ErrInternal       = errors.New("internal error")
)

type CalculatePointInput struct {
	MemberID       string
	PurchaseAmount int
}

type CalculatePointOutput struct {
	MemberID       string
	PurchaseAmount int
	GrantedPoint   int
}

type CalculatePointUsecase interface {
	Execute(ctx context.Context, input CalculatePointInput) (CalculatePointOutput, error)
}

type calculatePointUsecase struct {
	repository domain.MemberQueryRepository
}

func NewCalculatePointUsecase(repository domain.MemberQueryRepository) CalculatePointUsecase {
	return &calculatePointUsecase{
		repository: repository,
	}
}

func (u *calculatePointUsecase) Execute(
	ctx context.Context,
	input CalculatePointInput,
) (CalculatePointOutput, error) {
	member, err := u.repository.FindByID(ctx, nil, input.MemberID)
	if err != nil {
		if errors.Is(err, domain.ErrMemberNotFound) {
			return CalculatePointOutput{}, ErrMemberNotFound
		}
		return CalculatePointOutput{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	calculator := domain.NewPointCalculator()
	grantedPoint, err := calculator.Calculate(member, input.PurchaseAmount)
	if err != nil {
		return CalculatePointOutput{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return CalculatePointOutput{
		MemberID:       member.ID(),
		PurchaseAmount: input.PurchaseAmount,
		GrantedPoint:   grantedPoint,
	}, nil
}
