.PHONY: build test fmt vet lint run clean

build:
	go build -o bin/dirxray ./cmd/dirxray

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: vet
	@echo "Add golangci-lint locally: brew install golangci-lint && golangci-lint run"

run: build
	./bin/dirxray .

clean:
	rm -rf bin/
