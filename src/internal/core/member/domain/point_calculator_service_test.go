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
			member, err := NewMember("11111111-1111-1111-1111-111111111111", "田中", rank)
			assert.NoError(t, err)
			calculator := NewPointCalculator()

			point, err := calculator.Calculate(member, tt.purchaseAmount)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, point)
		})
	}

	t.Run("未知rankでエラーになること", func(t *testing.T) {
		member := &Member{
			id:   "11111111-1111-1111-1111-111111111111",
			name: "田中",
			rank: Rank("platinum"),
		}
		calculator := NewPointCalculator()

		_, err := calculator.Calculate(member, 5000)

		assert.Error(t, err)
	})
}
