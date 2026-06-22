package domain

import (
	"context"
)

type MemberRepository interface {
	FindByID(ctx context.Context, id string) (*Member, error)
}
