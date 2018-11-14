
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


build:
	GOOS=linux GOARCH=amd64 go build -o api-gateway main.go

test:
	go test -v ./... -cover

docker:
	docker build . -t api-gateway:latest
