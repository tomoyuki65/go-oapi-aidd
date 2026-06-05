.PHONY: generate

generate:
	docker compose run --rm api oapi-codegen -config oapi-codegen.models.yaml openapi/openapi.yaml
	docker compose run --rm api oapi-codegen -config oapi-codegen.server.yaml openapi/openapi.yaml