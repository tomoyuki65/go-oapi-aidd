package seed

import (
	"context"

	"github.com/uptrace/bun"
)

type Seeder interface {
	Seed(ctx context.Context, db *bun.DB) error
	Cleanup(ctx context.Context, db *bun.DB) error
}

func Run(ctx context.Context, db *bun.DB, seeders ...Seeder) error {
	for _, seeder := range seeders {
		if err := seeder.Seed(ctx, db); err != nil {
			return err
		}
	}
	return nil
}

func Cleanup(ctx context.Context, db *bun.DB, seeders ...Seeder) error {
	for i := len(seeders) - 1; i >= 0; i-- {
		if err := seeders[i].Cleanup(ctx, db); err != nil {
			return err
		}
	}
	return nil
}
