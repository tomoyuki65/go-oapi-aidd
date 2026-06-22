package usecase

import (
	"context"
	"errors"
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
