# 启动编译环境
FROM golang:1.16 as builder
WORKDIR /usr/src/app
#配置编译环境
RUN go env -w GOPROXY=https://goproxy.cn,direct
ARG TARGET_PATH

# 拷贝源代码
COPY ./go.mod ./
COPY ./go.sum ./
COPY . .
RUN CGO_ENABLED=0  go build -o server $TARGET_PATH/cmd

FROM alpine:3.12 as runner
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/server /opt/app/
CMD ["/opt/app/server"]

