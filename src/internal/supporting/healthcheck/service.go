package healthcheck

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	sl "go-oapi-aidd/internal/shared/logger"
)

type Service interface {
	Execute(ctx context.Context) error
}

type service struct {
	db     *bun.DB
	logger sl.Logger
}

func NewService(db *bun.DB, logger sl.Logger) Service {
	return &service{
		db:     db,
		logger: logger,
	}
}

func (s *service) Execute(ctx context.Context) error {
	// ログ出力
	s.logger.Info(true, ctx, "Healthcheck処理を実行（DB接続チェック含む）")

	// DB接続チェック
	if err := s.db.Ping(); err != nil {
		msg := fmt.Sprintf("failed to ping database: %s", err.Error())
		s.logger.Error(true, ctx, msg)
		return err
	}

	return nil
}
