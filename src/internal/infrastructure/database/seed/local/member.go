package local

import (
	"context"

	"github.com/uptrace/bun"

	"go-oapi-aidd/internal/infrastructure/database/schema"
)

const (
	TanakaMemberID = "11111111-1111-1111-1111-111111111111"
	SatoMemberID   = "22222222-2222-2222-2222-222222222222"
	SuzukiMemberID = "33333333-3333-3333-3333-333333333333"
)

type MemberSeeder struct{}

func NewMemberSeeder() MemberSeeder {
	return MemberSeeder{}
}

func (s MemberSeeder) Seed(ctx context.Context, db *bun.DB) error {
	members := []schema.Member{
		{ID: TanakaMemberID, Name: "田中", Rank: "bronze"},
		{ID: SatoMemberID, Name: "佐藤", Rank: "silver"},
		{ID: SuzukiMemberID, Name: "鈴木", Rank: "gold"},
	}
	_, err := db.NewInsert().Model(&members).On("CONFLICT (id) DO NOTHING").Exec(ctx)
	return err
}

func (s MemberSeeder) Cleanup(ctx context.Context, db *bun.DB) error {
	ids := []string{TanakaMemberID, SatoMemberID, SuzukiMemberID}
	_, err := db.NewDelete().Model((*schema.Member)(nil)).Where("id IN (?)", bun.List(ids)).Exec(ctx)
	return err
}
