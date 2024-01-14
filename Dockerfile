FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY cmd/      cmd/
COPY internal/ internal/
COPY go.sum    go.sum
COPY go.mod    go.mod

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates && \
    go build -o app cmd/app/main.go


FROM scratch AS runner

ENV GIN_MODE=release

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/app                           /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]
