unit_test:
	go test $(shell go list ./... | grep -v /tests)

integration_repo_test:
	go test -coverprofile=coverage.out -coverpkg ./database/postgres/repo/... ./tests/integration/repo

integration_service_test:
	go test -coverprofile=coverage.out -coverpkg ./internal/service/... ./tests/integration/service

coverage:
	go tool cover -html="coverage.out"

##swagger generate
swag_gen:
	swag init
