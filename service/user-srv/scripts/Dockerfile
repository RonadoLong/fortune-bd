FROM alpine

LABEL maintainer="Hong hongronado@gmail.com"
ARG configPath
ENV configPath ${configPath}
COPY user-srv user-srv
COPY conf.yaml conf.yaml
COPY confPro.yaml confPro.yaml

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
LABEL collectLog="true"
# 运行app
CMD "./user-srv" $configPath
