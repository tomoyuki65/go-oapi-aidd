package di

import (
	"github.com/uptrace/bun"

	memberDomain "go-oapi-aidd/internal/core/member/domain"
	memberQuery "go-oapi-aidd/internal/core/member/infrastructure/repository/query"
	memberUsecase "go-oapi-aidd/internal/core/member/usecase"
	sl "go-oapi-aidd/internal/shared/logger"
	"go-oapi-aidd/internal/supporting/healthcheck"
)

// DIコンテナの定義
type Container struct {
	DB                    *bun.DB
	Logger                sl.Logger
	HealthcheckService    healthcheck.Service
	MemberQueryRepository memberDomain.MemberQueryRepository
	CalculatePointUsecase memberUsecase.CalculatePointUsecase
}

// 依存関係の定義
type Dependencies struct {
	MemberQueryRepository memberDomain.MemberQueryRepository
	CalculatePointUsecase memberUsecase.CalculatePointUsecase
}

// 依存関係の上書き用関数のコンテナオプション定義
type ContainerOption func(*Dependencies)

// デフォルトの依存関係の作成関数
func NewDefaultDependencies() Dependencies {
	return Dependencies{}
}

// 依存関係からDIコンテナの作成関数
func NewContainerFromDependencies(db *bun.DB, logger sl.Logger, deps Dependencies) *Container {
	if deps.MemberQueryRepository == nil {
		deps.MemberQueryRepository = memberQuery.NewMemberRepository(db)
	}
	if deps.CalculatePointUsecase == nil {
		deps.CalculatePointUsecase = memberUsecase.NewCalculatePointUsecase(deps.MemberQueryRepository)
	}

	return &Container{
		DB:                    db,
		Logger:                logger,
		HealthcheckService:    healthcheck.NewService(db, logger),
		MemberQueryRepository: deps.MemberQueryRepository,
		CalculatePointUsecase: deps.CalculatePointUsecase,
	}
}

func WithMemberQueryRepository(repository memberDomain.MemberQueryRepository) ContainerOption {
	return func(deps *Dependencies) {
		deps.MemberQueryRepository = repository
	}
}

func WithCalculatePointUsecase(usecase memberUsecase.CalculatePointUsecase) ContainerOption {
	return func(deps *Dependencies) {
		deps.CalculatePointUsecase = usecase
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
