
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --gofast_out=. proto/*.proto

build:
	cd service/home-service && make build
	cd service/video-service && make build
	cd api-gateway && make build

test:
	go test -v ./... -cover

docker:
	docker build . -t home-service:latest
