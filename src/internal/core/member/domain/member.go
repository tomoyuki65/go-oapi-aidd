package domain

import "fmt"

type Member struct {
	id   string
	name string
	rank Rank
}

func NewMember(id string, name string, rank Rank) (*Member, error) {
	if !rank.IsValid() {
		return nil, fmt.Errorf("invalid rank: %s", rank)
	}
	return &Member{
		id:   id,
		name: name,
		rank: rank,
	}, nil
}

func (m *Member) ID() string {
	return m.id
}

func (m *Member) Name() string {
	return m.name
}

func (m *Member) Rank() Rank {
	return m.rank
}
