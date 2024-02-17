FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o card-auth-service ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/card-auth-service ./card-auth-service
ENTRYPOINT ./card-auth-service
