#!/bin/bash
set -xe
source ~/.bash_profile
find service -name "*.proto" | xargs -t -I{} protoc -I.:${GOPATH}/src --micro_out=. --gofast_out=plugins=micro:. {}
find service$1 -name "main.go"| xargs -t -I{} dirname shop-micro/{} | xargs -t -I{} go build {}