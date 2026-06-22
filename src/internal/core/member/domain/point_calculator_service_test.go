//go:build unit

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name           string
		rank           string
		purchaseAmount int
		want           int
	}{
		{name: "bronze 5,000円で50ptを返すこと", rank: "bronze", purchaseAmount: 5000, want: 50},
		{name: "silver 5,000円で150ptを返すこと", rank: "silver", purchaseAmount: 5000, want: 150},
		{name: "gold 15,000円で1,500ptを返すこと", rank: "gold", purchaseAmount: 15000, want: 1500},
		{name: "9,999円は倍率なしになること", rank: "gold", purchaseAmount: 9999, want: 499},
		{name: "10,000円は2倍になること", rank: "bronze", purchaseAmount: 10000, want: 200},
		{name: "小数点以下を切り捨てること", rank: "bronze", purchaseAmount: 999, want: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rank, err := NewRank(tt.rank)
			assert.NoError(t, err)
			amount, err := NewPurchaseAmount(tt.purchaseAmount)
			assert.NoError(t, err)

			point, err := NewPointCalculator().Calculate(rank, amount)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, point.Value())
		})
	}

	t.Run("未知rankでエラーになること", func(t *testing.T) {
		amount, err := NewPurchaseAmount(5000)
		assert.NoError(t, err)

		_, err = NewPointCalculator().Calculate(Rank("platinum"), amount)

		assert.Error(t, err)
	})
}

func TestNewPurchaseAmount(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{name: "0円を受け付けること", input: 0},
		{name: "999999999円を受け付けること", input: 999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, err := NewPurchaseAmount(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.input, amount.Value())
		})
	}

	testsError := []struct {
		name  string
		input int
	}{
		{name: "負数を拒否すること", input: -1},
		{name: "1000000000円を拒否すること", input: 1000000000},
	}

	for _, tt := range testsError {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPurchaseAmount(tt.input)

			assert.Error(t, err)
		})
	}
}

func TestNewGrantedPoint(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{name: "0ポイントを保持できること", input: 0},
		{name: "正のポイントを保持できること", input: 1500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point, err := NewGrantedPoint(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.input, point.Value())
		})
	}

	t.Run("負数ポイントを拒否すること", func(t *testing.T) {
		_, err := NewGrantedPoint(-1)

		assert.Error(t, err)
	})
}
