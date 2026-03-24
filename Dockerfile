# Этап 1: Сборка бинарного файла
FROM golang:1.22-alpine AS builder

# Устанавливаем необходимые инструменты (для gRPC и NATS может понадобиться)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник (статически скомпилированный для alpine)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/api .

# Этап 2: Финальный образ
FROM alpine:latest

# Устанавливаем CA-сертификаты для HTTPS-запросов
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из этапа сборки
COPY --from=builder /app/api .

# Открываем порт HTTP-сервера
EXPOSE 8080

# Запускаем сервис
CMD ["./api"]