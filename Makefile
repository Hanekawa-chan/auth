PROJECT_NAME := auth
PROJECT := github.com/Hanekawa-chan/auth

VERSION := $(shell git describe --tags)
COMMIT := $(shell git rev-parse --short HEAD)

LDFLAGS := "-s -w -X $(PROJECT)/internal/version.Version=$(VERSION) -X $(PROJECT)/internal/version.Commit=$(COMMIT)"
build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o ./bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)

test:
	@go test -v -cover -gcflags=-l --race ./...

GOLANGCI_LINT_VERSION := v1.50.0
lint:
	@golangci-lint run -v

dep:
	@go mod download

update_proto_win: update_proto proto_win generate_pb

update_proto_linux: update_proto proto_linux generate_pb

update_proto:
	git submodule update --remote

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
PWD := $(shell pwd)

generate_pb:
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) wontforgetit/protogen:2.0 -I=$(PROTO_PATH) --go_out=$(PROTO_PATH) --go-grpc_out=$(PROTO_PATH) --go-grpc_opt=require_unimplemented_servers=false `ls $(PROTO_PATH)`

update_deps:
	go get -u ./...

remove_tag:
	git tag -d "tag_name"

migrate:
	migrate create -dir migrations -ext sql auth

mock:
	mockgen -source=internal/app/interfaces.go -destination=internal/app/interfaces_mock_test.go -package=app -self_package=auth/internal/app