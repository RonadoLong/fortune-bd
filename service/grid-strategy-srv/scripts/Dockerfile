FROM alpine

ARG configPath
ENV configPath ${configPath}
COPY wq-grid-strategy wq-grid-strategy
COPY conf.toml conf.toml
COPY confPro.toml confPro.toml

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 9090
LABEL collectLog="true"
# 运行app
CMD "./wq-grid-strategy" $configPath