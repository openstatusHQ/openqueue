init:
    go install github.com/bufbuild/buf/cmd/buf@latest
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

build:
    go build -o ./tmp/main ./cmd/root.go

update:
    go mod tidy
    go get -u all
    go mod tidy

dev:
    air
