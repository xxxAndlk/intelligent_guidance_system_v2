# AI Diagnosis Service Dockerfile
# Multi-stage build for optimized production image

# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git make gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /ai-service ./cmd

# Runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /ai-service /app/ai-service
COPY --from=builder /build/configs/config.yaml /app/configs/config.yaml

ENV TZ=Asia/Shanghai
ENV SERVICE_NAME=ai-service
ENV SERVICE_PORT=8000

EXPOSE 8000 9000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/v1/health || exit 1

ENTRYPOINT ["/app/ai-service"]
CMD ["-config", "/app/configs/config.yaml"]