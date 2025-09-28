RAK BUKU RESTful API
Sebuah aplikasi RESTful API untuk manajemen rak buku digital yang dibangun dengan Go menggunakan framework Gin dan PostgreSQL. API ini menyediakan fitur untuk mengelola kategori buku dan koleksi buku dengan sistem autentikasi JWT.

🚀 Fitur
Autentikasi JWT - Login/Register user dengan token JWT
Manajemen Kategori - CRUD operasi untuk kategori buku
Manajemen Buku - CRUD operasi untuk buku dengan logic konversi ketebalan
Relasi Database - Kategori dan buku terhubung melalui foreign key
Dokumentasi Swagger - API documentation yang interaktif
Middleware JWT - Proteksi endpoint tertentu dengan autentikasi
Error Handling - Penanganan error yang konsisten
Database Migration - Sistem migrasi database yang terstruktur

🔐 Authentication
API ini menggunakan JWT (JSON Web Tokens) untuk autentikasi. Untuk mengakses endpoint yang memerlukan autentikasi:
Login melalui /api/users/login untuk mendapatkan token
Sertakan token dalam header Authorization: Bearer <your-jwt-token>

🎯 Business Logic
Konversi Ketebalan Buku
Sistem secara otomatis menentukan ketebalan buku berdasarkan jumlah halaman:
"tebal": Jika total_pages > 100
"tipis": Jika total_pages < 100

🧪 Testing dengan Postman
Import Collection: Import file Postman collection (jika ada)
Set Environment:

base_url: http://localhost:8080
jwt_token: (akan diisi setelah login)

Testing Flow:
Login untuk mendapatkan JWT token
Test CRUD operations untuk categories (dengan JWT)
Test CRUD operations untuk books
Test relasi categories dan books
