.PHONY: build test lint

build:
	go build -o conduit-connector-elasticsearch cmd/es/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

lint:
	golangci-lint run
