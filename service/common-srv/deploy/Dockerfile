# 启动编译环境
FROM golang:1.16 as builder

#配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码
WORKDIR /usr/src/app
COPY ./go.mod ./
COPY ./go.sum ./
COPY . .
RUN CGO_ENABLED=0  go build -o server ./service/common-srv/

FROM alpine:3.12 as runner
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/server /opt/app/
CMD ["/opt/app/server"]

