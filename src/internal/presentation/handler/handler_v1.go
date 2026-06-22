package handler

import (
	"context"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/presentation/gen"
	coreMember "go-oapi-aidd/internal/presentation/handler/core/member"
	"go-oapi-aidd/internal/presentation/handler/supporting/healthcheck"
)

type HandlerV1 struct {
	container *di.Container
}

func NewHandlerV1(
	container *di.Container,
) *HandlerV1 {
	return &HandlerV1{
		container: container,
	}
}

func (h *HandlerV1) Healthcheck(
	ctx context.Context,
	request gen.HealthcheckRequestObject,
) (gen.HealthcheckResponseObject, error) {
	healthcheckHandler := healthcheck.NewHealthcheckHandler(h.container, h.container.HealthcheckService)
	return healthcheckHandler.Healthcheck(ctx, request)
}

func (h *HandlerV1) CalculateMemberPoint(
	ctx context.Context,
	request gen.CalculateMemberPointRequestObject,
) (gen.CalculateMemberPointResponseObject, error) {
	memberHandler := coreMember.NewCalculatePointHandler(h.container.CalculatePointUsecase)
	return memberHandler.CalculateMemberPoint(ctx, request)
}
