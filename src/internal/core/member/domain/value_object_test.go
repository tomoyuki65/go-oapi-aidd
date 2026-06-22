//go:build unit

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRank(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "bronzeを受け付けること", input: "bronze"},
		{name: "silverを受け付けること", input: "silver"},
		{name: "goldを受け付けること", input: "gold"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rank, err := NewRank(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.input, rank.String())
		})
	}

	t.Run("未定義rankを拒否すること", func(t *testing.T) {
		_, err := NewRank("platinum")

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
