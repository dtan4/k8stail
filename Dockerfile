FROM alpine:3.4

RUN apk add --no-cache --update ca-certificates

COPY bin/k8stail /k8stail

ENTRYPOINT ["/k8stail"]
