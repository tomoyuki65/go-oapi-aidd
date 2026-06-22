//go:build integration

package member_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-oapi-aidd/internal/core/member/domain"
	"go-oapi-aidd/internal/core/member/infrastructure/repository/query"
	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/infrastructure/database/seed"
	"go-oapi-aidd/internal/infrastructure/database/seed/local"
)

func TestMemberRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	db, err := database.NewBunDB()
	assert.NoError(t, err)
	defer db.Close()

	memberSeeder := local.NewMemberSeeder()
	assert.NoError(t, seed.Run(ctx, db, memberSeeder))
	t.Cleanup(func() {
		assert.NoError(t, seed.Cleanup(ctx, db, memberSeeder))
	})

	repository := query.NewMemberRepository(db)

	t.Run("seed済み会員IDを使って指定IDの会員を取得できること", func(t *testing.T) {
		member, err := repository.FindByID(ctx, local.TanakaMemberID)

		assert.NoError(t, err)
		assert.Equal(t, local.TanakaMemberID, member.ID())
		assert.Equal(t, "田中", member.Name())
		assert.Equal(t, "bronze", member.Rank().String())
	})

	t.Run("存在しない会員IDがnot found扱いになること", func(t *testing.T) {
		_, err := repository.FindByID(ctx, "99999999-9999-9999-9999-999999999999")

		assert.ErrorIs(t, err, domain.ErrMemberNotFound)
	})
}
