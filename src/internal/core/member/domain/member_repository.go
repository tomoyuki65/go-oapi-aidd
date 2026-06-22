package domain

import (
	"context"

	"go-oapi-aidd/internal/shared/database"
)

type MemberRepository interface {
	FindByID(ctx context.Context, tx database.Transaction, id string) (*Member, error)
}
