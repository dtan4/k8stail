FROM golang:1.18 AS builder

WORKDIR /go/src/github.com/dtan4/k8stail

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /k8stail

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /k8stail /k8stail

ENTRYPOINT ["/k8stail"]
