.PHONY: build dev test lint clean install

BINARY := notespub
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.Version=$(VERSION)

build:  ## Compile CSS then build binary
	npx tailwindcss -i stylesheets/main.css -o style.css --minify
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/notespub

dev:    ## Watch mode: recompile on changes
	npx tailwindcss -i stylesheets/main.css -o style.css --watch &
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/notespub

clean:
	rm -f $(BINARY)

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/notespub

test:
	go test ./...

lint:
	go tool golangci-lint run
