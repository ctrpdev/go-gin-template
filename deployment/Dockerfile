# --- Etapa Base ---
FROM golang:1.25.0-alpine AS base

RUN apk add --no-cache git tzdata
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# --- Etapa Development (Live Reload) ---
FROM base AS dev
# Instalar Air para hot-reloading
RUN go install github.com/air-verse/air@latest
COPY . .
# Comando para levantar en desarrollo
CMD ["air", "-c", ".air.toml"]

# --- Etapa Builder (Compilación) ---
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o api_bin ./cmd/api/main.go

# --- Etapa Producción (Ligera) ---
FROM alpine:3.19 AS prod

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/api_bin .
COPY --from=builder /app/migrations ./migrations

ENV ENVIRONMENT=production
ENV SERVER_ADDRESS=:8080

EXPOSE 8080
CMD ["./api_bin"]
