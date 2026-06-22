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
