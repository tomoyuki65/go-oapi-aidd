package di

import (
	"github.com/uptrace/bun"

	sl "go-oapi-aidd/internal/shared/logger"
	"go-oapi-aidd/internal/supporting/healthcheck"
)

// DIコンテナの定義
type Container struct {
	DB                 *bun.DB
	Logger             sl.Logger
	HealthcheckService healthcheck.Service
}

// 依存関係の定義
type Dependencies struct{}

// 依存関係の上書き用関数のコンテナオプション定義
type ContainerOption func(*Dependencies)

// デフォルトの依存関係の作成関数
func NewDefaultDependencies() Dependencies {
	return Dependencies{}
}

// 依存関係からDIコンテナの作成関数
func NewContainerFromDependencies(db *bun.DB, logger sl.Logger, deps Dependencies) *Container {
	return &Container{
		DB:                 db,
		Logger:             logger,
		HealthcheckService: healthcheck.NewService(db, logger),
	}
}

// DIコンテナの作成関数
func NewContainer(db *bun.DB, logger sl.Logger, opts ...ContainerOption) *Container {
	// デフォルトの依存関係の取得
	deps := NewDefaultDependencies()

	// 依存関係の上書き処理
	for _, opt := range opts {
		opt(&deps)
	}

	return NewContainerFromDependencies(db, logger, deps)
}
