package healthcheck

import (
	"context"

	"go-oapi-aidd/internal/di"
	"go-oapi-aidd/internal/presentation/gen"
	"go-oapi-aidd/internal/supporting/healthcheck"
)

type HealthcheckHandler struct {
	container          *di.Container
	healthcheckService healthcheck.Service
}

func NewHealthcheckHandler(container *di.Container, healthcheckService healthcheck.Service) *HealthcheckHandler {
	return &HealthcheckHandler{
		container:          container,
		healthcheckService: healthcheckService,
	}
}

func (h *HealthcheckHandler) Healthcheck(
	ctx context.Context,
	request gen.HealthcheckRequestObject,
) (gen.HealthcheckResponseObject, error) {
	if err := h.healthcheckService.Execute(ctx); err != nil {
		return gen.Healthcheck500JSONResponse{
			Message: "Internal Server Error",
		}, nil
	} else {
		return gen.Healthcheck200JSONResponse{
			Message: "OK",
		}, nil
	}
}
