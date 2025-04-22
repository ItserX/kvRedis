FROM golang:1.24.1

WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache wait-for-it

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY cmd ./cmd
COPY internal ./internal
COPY .env .

# Собираем приложение
RUN go build -o my-app ./cmd/main.go

# Команда запуска с ожиданием готовности Redis
CMD ["sh", "-c", "wait-for-it redis:6379 --timeout=30 --strict -- ./my-app"]