# TokoLoka - API Gateway

TokoLoka adalah **API Gateway** yang dirancang untuk mengelola data dan laporan dari berbagai layanan backend. Dengan fitur autentikasi berbasis **JWT**, kontrol akses berbasis peran (**RBAC**), dan kemampuan pengelolaan data serta laporan, TokoLoka mendukung administrator dan pengguna biasa untuk mengakses sistem dengan aman dan efisien.

---

## Fitur Utama

- **Manajemen Data**
    - **Kategori**: Tambah, ubah, hapus, dan lihat kategori.
    - **Produk**: Tambah, ubah, hapus, lihat produk, serta unggah gambar produk.
    - **Transaksi**: Lihat, buat, dan kelola transaksi dengan dukungan nomor tujuan (**destination number**) dan nomor seri (**serial number**).

- **Pelaporan**
    - Membuat laporan berdasarkan filter (tanggal, pengguna, produk, dll.).
    - Laporan dapat diunduh dalam format **CSV** atau **PDF**.
    - Mendukung paginasi untuk data besar.

- **Logging Aktivitas**
    - Mencatat semua aktivitas penting, seperti login, pembuatan laporan, dan perubahan data.

- **Kontrol Akses Berbasis Peran (RBAC)**
    - **Administrator** memiliki akses penuh untuk mengelola semua data.
    - **Pengguna biasa** memiliki akses terbatas sesuai kebutuhan.

---

## Teknologi

- **Backend Framework**: [Gin](https://gin-gonic.com/) - Go web framework
- **Database**: MySQL, menggunakan ORM GORM
- **Autentikasi**: JSON Web Token (JWT)
- **Logging**: Uber Zap
- **Penyimpanan File**: Lokal untuk unggah gambar

---

## Instalasi

### Prasyarat
- **Go** (minimal versi 1.19)
- **MySQL**
- **Postman** atau alat pengujian API lainnya

### Langkah-Langkah
1. Clone repository ini:
   ```bash
   git clone https://github.com/your-username/tokoloka.git
   cd tokloka

### Buat file .env dan isi seperti berikut:
- DB_HOST=localhost
- DB_PORT=3306/3307/etc
- DB_USER=root
- DB_PASS=password
- DB_NAME=tokoloka
- JWT_SECRET=your_jwt_secret

### Jalankan perintah untuk menginstal dependensi:
go mod tidy

### Jalankan aplikasi:
go run main.go

### Akses API Gateway melalui http://localhost:8080.

## Dokumentasi API
### Autentikasi
- POST /auth/register - Registrasi pengguna baru
- POST /auth/login - Login pengguna
### Manajemen Kategori
- POST /api/categories - Tambah kategori
- PUT /api/categories/:id - Ubah kategori
- DELETE /api/categories/:id - Hapus kategori
- GET /api/categories - Lihat semua kategori
### Manajemen Produk
- POST /api/products - Tambah produk
- PUT /api/products/:id - Ubah produk
- DELETE /api/products/:id - Hapus produk
- POST /api/products/:id/image - Unggah gambar produk
- GET /api/products - Lihat semua produk
- GET /api/products/:id - Lihat detail produk
### Manajemen Transaksi
- POST /api/transactions - Buat transaksi baru
- GET /api/transactions - Lihat semua transaksi
- GET /api/transactions/:id - Lihat detail transaksi
### Laporan
- POST /api/reports/generate - Membuat laporan berdasarkan filter
- GET /api/reports/download - Mengunduh laporan dalam format CSV atau PDF

## Struktur Proyek
TokoLoka/
├── main.go                     # Entry point aplikasi
├── config/                     # Konfigurasi database
├── controller/                 # HTTP handlers dan logika endpoint
├── entity/                     # Definisi model database
├── middleware/                 # Middleware untuk autentikasi dan logging
├── repository/                 # Akses database
├── service/                    # Logika bisnis
├── uploads/                    # Direktori untuk file gambar produk
└── logs/                       # Direktori untuk file log aktivitas

## Use Case
### Administrator
Mengelola kategori, produk, dan laporan.
Mengakses semua data transaksi dan laporan.
Melihat aktivitas sistem melalui log.
### User
Melihat daftar produk dan kategori.
Melakukan transaksi dan melihat riwayatnya.
Mengunduh laporan transaksi mereka sendiri.

## Lisensi
TokoLoka dilisensikan di bawah MIT License.
