FROM gcr.io/distroless/static

COPY k8stail /

ENTRYPOINT ["/k8stail"]
