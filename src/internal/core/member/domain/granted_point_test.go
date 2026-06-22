//go:build unit

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
