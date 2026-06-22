package domain

import "fmt"

type Rank string

const (
	RankBronze Rank = "bronze"
	RankSilver Rank = "silver"
	RankGold   Rank = "gold"
)

func NewRank(value string) (Rank, error) {
	rank := Rank(value)
	if !rank.IsValid() {
		return "", fmt.Errorf("invalid rank: %s", value)
	}
	return rank, nil
}

func (r Rank) String() string {
	return string(r)
}

func (r Rank) IsValid() bool {
	switch r {
	case RankBronze, RankSilver, RankGold:
		return true
	default:
		return false
	}
}

func (r Rank) pointRatePercent() (int, error) {
	switch r {
	case RankBronze:
		return 1, nil
	case RankSilver:
		return 3, nil
	case RankGold:
		return 5, nil
	default:
		return 0, fmt.Errorf("invalid rank: %s", r)
	}
}
