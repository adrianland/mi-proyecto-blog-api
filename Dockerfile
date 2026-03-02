# Stage 1: Builder
FROM golang:1.26 AS builder

WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar aplicación (estática)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o api ./cmd/api

# Stage 2: Runtime
FROM debian:bookworm-slim

# Instalar solo lo necesario y crear usuario en un solo paso para reducir capas
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -g 1000 appgroup && useradd -u 1000 -g appgroup appuser

WORKDIR /app

# Copiar solo el binario desde builder (sin el .env)
COPY --from=builder /app/api .

# Cambiar permisos y asegurar que el usuario tenga acceso
RUN chown -R appuser:appgroup /app

# Cambiar al usuario no-root
USER appuser

# Exponer el puerto (este valor es solo informativo, Docker Compose manda)
EXPOSE 8080

# Comando de inicio
CMD ["./api"]