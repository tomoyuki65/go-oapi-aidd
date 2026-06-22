package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	"go-oapi-aidd/internal/core/member/domain"
	"go-oapi-aidd/internal/infrastructure/database/schema"
)

type MemberRepository struct {
	db bun.IDB
}

func NewMemberRepository(db ...bun.IDB) *MemberRepository {
	var executor bun.IDB
	if len(db) > 0 {
		executor = db[0]
	}
	return &MemberRepository{
		db: executor,
	}
}

func (r *MemberRepository) FindByID(ctx context.Context, db bun.IDB, id string) (*domain.Member, error) {
	if db == nil {
		db = r.db
	}
	if db == nil {
		return nil, errors.New("db is nil")
	}

	member := new(schema.Member)
	err := db.NewSelect().
		Model(member).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMemberNotFound
		}
		return nil, err
	}

	rank, err := domain.NewRank(member.Rank)
	if err != nil {
		return nil, fmt.Errorf("convert member rank: %w", err)
	}

	return domain.NewMember(member.ID, member.Name, rank)
}
