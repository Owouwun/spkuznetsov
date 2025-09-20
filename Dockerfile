FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /app/main ./cmd/server

FROM scratch

COPY --from=builder /app/main /main
COPY --from=builder /app/migrations /migrations

EXPOSE 8080

ENTRYPOINT ["/main"]