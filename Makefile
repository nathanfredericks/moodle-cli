BINARY := moodle
VERSION ?= dev
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: build build-lambda test lint clean generate docs sam-validate sam-build

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/moodle

build-lambda:
	docker build --platform linux/arm64 -f Dockerfile.lambda --build-arg VERSION=$(VERSION) -t moodle-mcp:local .

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
	go run ./tools/gendocs -out skills/moodle-cli/references

sam-validate:
	sam validate --lint --template-file template.yaml
	sam validate --lint --template-file infrastructure/bootstrap.yaml

sam-build:
	sam build --template-file template.yaml

install: build
	mv $(BINARY) $(GOPATH)/bin/
