package domain

import (
	"context"
	"errors"
)

var ErrMemberNotFound = errors.New("member not found")

type Member struct {
	id   string
	name string
	rank Rank
}

type Rank string

type MemberRepository interface {
	FindByID(ctx context.Context, id string) (*Member, error)
}
