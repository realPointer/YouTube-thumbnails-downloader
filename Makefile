LOCAL_BIN:=$(CURDIR)/bin

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume:
	docker volume rm youtube-thumbnails-downloader_thumbnails

install-all: ### Install all dependencies and tools
	make install-deps
	make install-linter
.PHONY: install-all

install-deps: ### Install dependencies
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
.PHONY: install-deps

generate: ### Generate all proto files
	make generate-thumbnail-api
.PHONY: generate

generate-thumbnail-api: ### Generate thumbnail api
	mkdir -p pkg/thumbnail_v1
	protoc --proto_path api/thumbnail_v1 \
	--go_out=pkg/thumbnail_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/thumbnail_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/thumbnail_v1/thumbnail.proto		
.PHONY: generate-thumbnail-api

install-linter: ### Install golangci-lint
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
.PHONY: install-linter

lint: ### Check by golangci linter
	$(LOCAL_BIN)/golangci-lint run
.PHONY: lint

test: ### Run test
	go test -v ./...
.PHONY: test