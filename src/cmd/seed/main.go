package main

import (
	"context"
	"log/slog"
	"os"

	"go-oapi-aidd/internal/infrastructure/database"
	"go-oapi-aidd/internal/infrastructure/database/seed"
	"go-oapi-aidd/internal/infrastructure/database/seed/local"
)

func main() {
	db, err := database.NewBunDB()
	if err != nil {
		slog.Error("failed to connect database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := seed.Run(context.Background(), db, local.NewMemberSeeder()); err != nil {
		slog.Error("failed to run seed", "err", err)
		os.Exit(1)
	}
}
