
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker

proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --gofast_out=. proto/*.proto

build: proto
	GOOS=linux GOARCH=amd64 go build -o info-service main.go plugin.go

test:
	go test -v ./... -cover

docker:
	docker build . -t info-service:latest
