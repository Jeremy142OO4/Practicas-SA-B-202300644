FROM golang:1.23-alpine AS compilacion
WORKDIR /codigo
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /salida/api ./cmd/api

FROM alpine:3.21
RUN apk add --no-cache ca-certificates wget && addgroup -S aplicacion && adduser -S aplicacion -G aplicacion
WORKDIR /app
COPY --from=compilacion /salida/api ./api
USER aplicacion
EXPOSE 8080
CMD ["./api"]
