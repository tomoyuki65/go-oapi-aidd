//go:build unit

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
