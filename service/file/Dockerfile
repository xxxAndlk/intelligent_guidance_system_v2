FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-service ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/file-service .
COPY --from=builder /app/configs ./configs

ENV TZ=Asia/Shanghai

EXPOSE 8012 9012

CMD ["./file-service"]