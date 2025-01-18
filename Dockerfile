# Gunakan image Golang sebagai base image
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Salin go.mod dan go.sum untuk menginstall dependencies
COPY go.mod go.sum ./

# Install dependencies
RUN go mod tidy

# Salin seluruh kode sumber ke dalam container
COPY . .

# Build aplikasi
RUN go build -o main .

# Expose port yang digunakan oleh aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
