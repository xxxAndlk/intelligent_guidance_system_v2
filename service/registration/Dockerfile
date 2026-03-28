FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o registration-service ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/registration-service .
COPY --from=builder /app/configs ./configs

ENV TZ=Asia/Shanghai

EXPOSE 8007 9007

CMD ["./registration-service"]