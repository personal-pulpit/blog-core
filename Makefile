## unit test
unit_test:
	go test $(shell go list ./... | grep -v /tests)
