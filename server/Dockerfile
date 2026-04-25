# Stage 1: Build (Gunakan versi yang lebih spesifik jika memungkinkan)
FROM golang:alpine AS builder

# Instal git dan build-base (terkadang diperlukan untuk beberapa module Go)
RUN apk add --no-cache git

WORKDIR /app

# Strategi Caching: Copy file modul dulu agar tidak download ulang tiap kali kode berubah
COPY go.mod go.sum ./
RUN go mod download

# Baru copy semua source code
COPY . .

# Build dengan flag tambahan agar binary statis (aman di Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Deployment (Gunakan Alpine yang ringan)
FROM alpine:latest  

# Tambahkan sertifikat SSL agar backend bisa melakukan request HTTPS ke API Gemini
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Salin binary dari stage builder
COPY --from=builder /app/main .

# Pastikan file binary punya akses eksekusi
RUN chmod +x ./main

# Port backend (sesuaikan dengan main.go kamu)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]