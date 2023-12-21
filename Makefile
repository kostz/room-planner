TARGET ?= room-planner

build: build/$(TARGET)

build/$(TARGET): $(SOURCES)
	go build -o build/ -ldflags "-s"

test:
	go test -race -timeout=10s ./...

lint:
	go version
	golangci-lint run
	go mod tidy

fmt:
	go fmt ./...
	go mod tidy

clean:
	@rm -f build/$(TARGET)
