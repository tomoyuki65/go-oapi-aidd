package domain

import "errors"

const pointMultiplierThreshold = 10000

var ErrInvalidPurchaseAmount = errors.New("invalid purchase amount")

type PointCalculator struct{}

func NewPointCalculator() *PointCalculator {
	return &PointCalculator{}
}

func (c *PointCalculator) Calculate(member *Member, purchaseAmount int) (int, error) {
	if member == nil {
		return 0, ErrMemberNotFound
	}
	if purchaseAmount < 0 {
		return 0, ErrInvalidPurchaseAmount
	}

	rate, err := member.Rank().pointRatePercent()
	if err != nil {
		return 0, err
	}

	point := purchaseAmount * rate / 100
	if purchaseAmount >= pointMultiplierThreshold {
		point *= 2
	}

	return point, nil
}
