FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o profile-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/profile-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./profile-service"] 