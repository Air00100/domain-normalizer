.PHONY: all fa fmt lint test cover coverfn

all: fa fmt lint test

fa:
	@fieldalignment -fix ./...

fmt:
	@goimports -w -local normalizer .
	@gofmt -w .
	@golines -w .

lint:
	@golangci-lint run

test:
	@go test ./...

cover:
	go test -coverpkg=./... -coverprofile=coverage.out ./... && go tool $@ -html=coverage.out

coverfn:
	go test -coverpkg=./... -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out
