FROM ubuntu:20.04
RUN mkdir work
WORKDIR /work
COPY ./go1.14.2.linux-amd64.tar.gz ./
RUN tar -C /usr/local -zxf go1.14.2.linux-amd64.tar.gz
RUN export GOROOT=/usr/local/go && export PATH=$PATH:$GOROOT/bin
RUN /usr/local/go/bin/go version
RUN apt install git
RUN /usr/local/go/bin/go get -v -u github.com/micro/micro/v2
RUN echo micro version
RUN export MICRO_REGISTRY=etcd && export MICRO_REGISTRY_ADDRESS=192.168.101.192:2379
CMD [ "micro web" ]