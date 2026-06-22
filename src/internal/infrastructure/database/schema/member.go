package schema

import "time"

type Member struct {
	ID        string    `bun:"id,pk,type:uuid"`
	Name      string    `bun:"name,notnull"`
	Rank      string    `bun:"rank,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
}
