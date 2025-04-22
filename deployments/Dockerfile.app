FROM golang:1.24.1

WORKDIR /app

RUN apk add --no-cache wait-for-it

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY .env .

RUN go build -o my-app ./cmd/main.go

CMD ["sh", "-c", "wait-for-it redis:6379 --timeout=30 --strict -- ./my-app"]
