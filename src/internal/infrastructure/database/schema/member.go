package schema

import (
	"time"

	"github.com/uptrace/bun"
)

type Member struct {
	bun.BaseModel `bun:"table:members,alias:m"`

	ID        string    `bun:"id,pk,type:uuid"`
	Name      string    `bun:"name,notnull"`
	Rank      string    `bun:"rank,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
}
