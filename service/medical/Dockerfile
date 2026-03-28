FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/medical-service ./cmd/

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/medical-service .
COPY --from=builder /app/configs ./configs

EXPOSE 8004 9004

ENTRYPOINT ["./medical-service"]
CMD ["-conf", "./configs/config.yaml"]