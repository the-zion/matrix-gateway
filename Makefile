API_PROTO_FILES=$(shell find api -name *.proto)
VERSION=$(shell git describe --tags --abbrev=0)
.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
           --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: tidy
# generate
tidy:
	go mod tidy

.PHONY: all
# generate all
all:
	make init;
	make api;
	make tidy;
	make build;

.PHONY: image
# image
image:
	make all;
	docker build -t gateway:$(VERSION) .
