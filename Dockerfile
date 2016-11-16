FROM alpine:3.4

ENV GLIBC_VERSION 2.23-r3

RUN apk add --no-cache --update ca-certificates

RUN apk add --no-cache --update wget \
    && wget -qO /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub \
    && wget -q https://github.com/sgerrand/alpine-pkg-glibc/releases/download/$GLIBC_VERSION/glibc-$GLIBC_VERSION.apk \
    && apk add glibc-$GLIBC_VERSION.apk \
    && rm glibc-$GLIBC_VERSION.apk \
    && apk del wget

COPY bin/k8stail /k8stail

ENTRYPOINT ["/k8stail"]
