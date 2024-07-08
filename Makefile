## unit test
unit_test:
	go test -v $(shell go list ./... | grep -v /tests)