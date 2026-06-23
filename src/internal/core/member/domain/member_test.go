//go:build unit

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMember(t *testing.T) {
	t.Run("id、name、rankを保持できること", func(t *testing.T) {
		rank, err := NewRank("gold")
		assert.NoError(t, err)

		member, err := NewMember("11111111-1111-1111-1111-111111111111", "田中", rank)

		assert.NoError(t, err)
		assert.Equal(t, "11111111-1111-1111-1111-111111111111", member.ID())
		assert.Equal(t, "田中", member.Name())
		assert.Equal(t, rank, member.Rank())
	})

	t.Run("不正なrankを拒否すること", func(t *testing.T) {
		_, err := NewMember("11111111-1111-1111-1111-111111111111", "田中", Rank("platinum"))

		assert.Error(t, err)
	})
}
