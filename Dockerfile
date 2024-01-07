FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY cmd/      cmd/
COPY internal/ internal/
COPY go.sum    go.sum
COPY go.mod    go.mod

RUN go build -o app cmd/openapi/main.go


FROM scratch AS runner

ARG ENVIRONMENT
ARG ALLOW_ORIGIN
ARG TRUSTED_PROXY

ENV GIN_MODE=release \
    ENVIRONMENT=$ENVIRONMENT \
    ALLOW_ORIGIN=$ALLOW_ORIGIN \
    TRUSTED_PROXY=$TRUSTED_PROXY

COPY --from=builder /app/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]
