# Étape 1 : compilation du projet Go
FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Nécessaire pour github.com/mattn/go-sqlite3
RUN apt-get update && apt-get install -y gcc libc6-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./back-end/cmd


# Étape 2 : image finale plus légère
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/forum /app/forum

COPY front-end /app/front-end
COPY back-end/database /app/back-end/database

EXPOSE 8080

CMD ["/app/forum"]