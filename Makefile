## unit test
unit_test:
	go test $(shell go list ./... | grep -v /tests)

integration_test:
	go test -cover -coverpkg ./internal/repository/... ./tests/integration/repo
