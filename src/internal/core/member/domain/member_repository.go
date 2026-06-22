package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type MemberRepository interface {
	FindByID(ctx context.Context, db bun.IDB, id string) (*Member, error)
}
