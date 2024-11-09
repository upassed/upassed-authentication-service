FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o upassed-authentication-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir -p /upassed-authentication-service/config
RUN mkdir -p /upassed-authentication-service/migration/scripts
COPY --from=builder /app/upassed-authentication-service /upassed-authentication-service/upassed-authentication-service
COPY --from=builder /app/config/* /upassed-authentication-service/config
COPY --from=builder /app/migration/scripts/* /upassed-authentication-service/migration/scripts
RUN chmod +x /upassed-authentication-service/upassed-authentication-service
ENV APP_CONFIG_PATH="/upassed-authentication-service/config/local.yml"
EXPOSE 44045
CMD ["/upassed-authentication-service/upassed-authentication-service"]
