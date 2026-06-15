.PHONY: generate \
	test test-unit test-integration test-e2e

generate:
	docker compose run --rm api oapi-codegen -config oapi-codegen.models.yaml openapi/openapi.yaml
	docker compose run --rm api oapi-codegen -config oapi-codegen.server.yaml openapi/openapi.yaml

test: test-unit test-integration test-e2e

test-unit:
	docker compose exec api env ENV=testing go test -v -cover -tags=unit $$(docker compose exec api env ENV=testing go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' -tags=unit ./...)

test-integration:
	docker compose exec api env ENV=testing go test -v -tags=integration $$(docker compose exec api env ENV=testing go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' -tags=integration ./...)

test-e2e:
	docker compose exec api env ENV=testing go test -v -tags=e2e $$(docker compose exec api env ENV=testing go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' -tags=e2e ./...)
