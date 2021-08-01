FROM alpine
LABEL maintainer="Hong hongronado@gmail.com"
# 修改时区
# RUN apk add -U tzdata
# RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 复制app文件到工作目录下
ARG configPath
ENV configPath ${configPath}
COPY app app
COPY conf.yaml conf.yaml
COPY confPro.yaml confPro.yaml

LABEL collectLog="true"
# 运行app
CMD "./app" $configPath
