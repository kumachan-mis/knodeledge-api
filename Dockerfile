FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY cmd/      cmd/
COPY internal/ internal/
COPY go.sum    go.sum
COPY go.mod    go.mod

RUN go build -o app cmd/app/main.go


FROM scratch AS runner

COPY --from=builder /app/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]
