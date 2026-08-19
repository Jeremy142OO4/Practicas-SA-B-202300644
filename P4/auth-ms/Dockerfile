FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-ms .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-ms .

EXPOSE 4001

CMD ["./auth-ms"]
