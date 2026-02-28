BINARY := moodle
VERSION ?= dev
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: build test lint clean generate docs

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/moodle

test:
	go test ./... -v -count=1

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/

generate:
	go generate ./...

docs:
	go run ./tools/gendocs -out docs/

install: build
	mv $(BINARY) $(GOPATH)/bin/
