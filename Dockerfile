FROM alpine:3.4

ENV GLIBC_VERSION 2.23-r3

RUN apk add --no-cache --update ca-certificates

COPY bin/k8stail /k8stail

ENTRYPOINT ["/k8stail"]
