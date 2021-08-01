
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/forwardOffer.proto
	protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/exorder.proto

.PHONY: build
build: proto
	go build -o forward-offer-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t forward-offer-srv:latest
