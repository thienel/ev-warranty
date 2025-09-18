FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# ----------- final image -----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

# Copy binary
COPY --from=builder /app/main .

# Copy entrypoint
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

# Entrypoint sẽ chạy main và tạo keys từ env
ENTRYPOINT ["./entrypoint.sh"]
