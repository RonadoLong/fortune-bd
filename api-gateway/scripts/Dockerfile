FROM alpine:latest
ARG configPath
ENV configPath ${configPath}
RUN apk --no-cache add ca-certificates
COPY api-gateway api-gateway
COPY conf.yaml conf.yaml
COPY confPro.yaml confPro.yaml

LABEL collectLog="true"
EXPOSE 9530
CMD "./api-gateway" $configPath
