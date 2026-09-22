FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiamos todo el código fuente primero
COPY . .

# Descargamos dependencias y generamos el go.sum automáticamente
RUN go mod tidy

# Compilamos el binario de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Imagen final ultra liviana
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]