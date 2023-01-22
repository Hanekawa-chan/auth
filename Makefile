PROJECT_NAME := auth
PROJECT := github.com/kanji-team/auth

VERSION := $(shell git describe --tags)
COMMIT := $(shell git rev-parse --short HEAD)

LDFLAGS := "-s -w -X $(PROJECT)/internal/version.Version=$(VERSION) -X $(PROJECT)/internal/version.Commit=$(COMMIT)"
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o ./bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)

test:
	@go test -v -cover -gcflags=-l --race ./...

GOLANGCI_LINT_VERSION := v1.24.0
lint:
	@golangci-lint run -v

dep:
	@go mod download

update_proto_win: update_proto proto_win generate_pb

update_proto_linux: update_proto proto_linux generate_pb

update_proto:
	git submodule foreach git pull origin main

proto_linux:
	rm -rf ./proto/protocol
	mkdir -p ./proto/protocol
	cp ./submodule/protocol/* ./proto/protocol # хз копирует ли оно папки
	rm -rf ./proto/services/*.pb.go
	mkdir -p proto/services

proto_win:
	rmdir /s /q proto\protocol
	mkdir proto\protocol
	robocopy submodule\protocol proto\protocol /MIR
	rmdir /s /q proto\services
	mkdir proto\services

PROTO_PATH := "proto/protocol"

generate_pb:
	docker run --rm -v $(pwd):$(pwd) -w $(pwd) protogen -I=$(PROTO_PATH) --go_out=$(PROTO_PATH) --go-grpc_out=$(PROTO_PATH) `ls $(PROTO_PATH)`

update_deps:
	go get -u ./...

remove_tag:
	git tag -d "tag_name"