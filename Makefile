## unit test
unit_test:
	go test $(shell go list ./... | grep -v /tests)

integration_test:
	go test -coverprofile=coverage.out -coverpkg ./database/postgres/repo/... ./tests/integration/repo
coverage:
	go tool cover -html="coverage.out"

