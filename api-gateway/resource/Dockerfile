FROM scratch
ADD main /main
ADD resource/pri_key.pem resource/pub_key.pem /resource/
ADD resource/dev.yaml resource/pro.yaml /resource/
ADD resource/ca-certificates.crt /etc/ssl/certs/
VOLUME /resource/
EXPOSE 9527
ENTRYPOINT ["/main"]