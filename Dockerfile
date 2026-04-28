# ========================
# 1. Build stage
# ========================
# Gunakan versi 1.25.4 sesuai permintaan go.mod
FROM golang:1.25.4-alpine AS builder

# Install git karena go mod mungkin membutuhkannya untuk private repo/dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files dulu agar layer cache tidak pecah jika hanya kode yang berubah
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary secara statis agar bisa jalan di distroless
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o clm-be ./main.go

# ========================
# 2. Runtime stage
# ========================
FROM alpine:3.23

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/go-archetype /app/app

# Copy folder migrations (pastikan folder ini ada di root project kamu)
COPY --from=builder /app/migrations/ /app/migrations/

# Copy dokumentasi Swagger
COPY --from=builder /app/internal/adapters/http/docs /app/internal/adapters/http/docs

# Expose port sesuai aplikasi
EXPOSE 8080

# ENTRYPOINT adalah executable utama
ENTRYPOINT ["/app/app"]

# CMD memberikan argumen default. 
# Ini akan menjalankan perintah: /app/app http
# Jika di K8s kamu menentukan args: ["migrate"], maka "http" ini akan otomatis diganti.
CMD ["http"]

# To use on k8s, add this on 'containers' field
# command: ["/app/app"]
# args: ["http"] or ["consumer"] or ["migrate", "up]
