# Etapa 1: Build
FROM golang:1.24-alpine AS builder

# Instalar dependencias necesarias para compilar y para herramientas de red o CGO (si hiciera falta)
RUN apk add --no-cache git tzdata

WORKDIR /app

# Primero copiar go.mod y go.sum para aprovechar la cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar la aplicación removiendo metadata de debug para reducir tamaño (-s -w)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o api_bin ./cmd/api/main.go

# Etapa 2: Imagen final ligéra
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copiar el binario
COPY --from=builder /app/api_bin .

# Variables de entorno por defecto (se pueden sobreescribir desde docker-compose)
ENV ENVIRONMENT=production
ENV SERVER_ADDRESS=:8080

EXPOSE 8080

# Ejecutar el binario
CMD ["./api_bin"]