# Finance App API 💰

API sederhana untuk mengelola keuangan pribadi Anda. Catat pemasukan dan pengeluaran, pantau saldo, dan analisis pengeluaran Anda dengan mudah.

## Fitur Utama 🚀

* **Manajemen Transaksi:**
    * Tambahkan, edit, dan hapus transaksi (pemasukan dan pengeluaran).
    * Kategorikan transaksi untuk analisis yang lebih baik.
    * Filter transaksi berdasarkan rentang tanggal.
* **Manajemen Kategori:**
    * Tambahkan, edit, dan hapus kategori pemasukan dan pengeluaran.
* **Saldo Real-time:**
    * Lihat saldo Anda saat ini secara real-time.
* **Laporan Keuangan:**
    * Lihat ringkasan pengeluaran dan pemasukan bulanan.
    * Ekspor laporan keuangan dalam format CSV atau PDF.
* **Autentikasi & Otorisasi:**
    * Sistem login dengan JWT (JSON Web Tokens) untuk menjaga keamanan data.
* **API Dokumentasi:**
    * Dokumentasi otomatis menggunakan Swagger/OpenAPI.

## Teknologi yang Digunakan 💻

* **Backend:** Go (Golang) dengan pendekatan *native* (tanpa framework web tambahan)
* **Database:** MongoDB (NoSQL database yang fleksibel dan mudah digunakan)
* **Autentikasi:** JSON Web Token (JWT) untuk autentikasi API

## Instalasi dan Penggunaan 🛠️

### Prasyarat

* **Go:** Pastikan Anda telah menginstal Go di sistem Anda. Anda dapat mengunduhnya dari [https://golang.org/](https://golang.org/).
* **MongoDB:** Pastikan Anda memiliki MongoDB yang berjalan di lokal atau di cloud. Anda dapat mengunduhnya dari [https://www.mongodb.com/](https://www.mongodb.com/).
* **Git:** Instal Git untuk mengkloning repositori ini. Dapat diunduh dari [https://git-scm.com/](https://git-scm.com/).

### Langkah-langkah

1. **Clone Repositori:**
    ```bash
    git clone https://github.com/username/finance_app.git
    cd finance_app
    ```

2. **Atur Environment Variables:**
    Buat file `.env` di root proyek dan tambahkan konfigurasi berikut:
    ```env
    MONGODB_URI=mongodb://localhost:27017/finance_app
    JWT_SECRET=your_jwt_secret
    PORT=8080
    ```

3. **Instal Dependencies:**
    Jalankan perintah berikut untuk menginstal dependensi yang diperlukan:
    ```bash
    go mod tidy
    ```

4. **Jalankan Aplikasi:**
    Setelah semua dependensi diinstal, jalankan aplikasi dengan perintah:
    ```bash
    go run main.go
    ```

5. **Akses API:**
    API akan berjalan di `http://localhost:8080`. Anda dapat menggunakan alat seperti Postman untuk mengakses endpoint API.

## Struktur Proyek 📂

Berikut adalah struktur direktori proyek ini:

```bash
finance_app/
├── config/             # Konfigurasi aplikasi
├── controllers/        # Logika bisnis dan pengendali HTTP
├── models/             # Struktur data dan model database
├── routes/             # Definisi rute API
├── services/           # Layanan untuk logika bisnis
├── utils/              # Fungsi utilitas dan helper
├── main.go             # Entry point aplikasi
└── .env.example        # Contoh file environment variables

## Dokumentasi API 📄

API ini mendukung dokumentasi otomatis menggunakan Swagger/OpenAPI. Setelah aplikasi berjalan, Anda dapat mengakses dokumentasi API di `http://localhost:8080/swagger/index.html`.

### Endpoints Utama

1. **Manajemen Transaksi:**
    - **POST** `/transactions`: Menambahkan transaksi baru.
    - **GET** `/transactions`: Mendapatkan daftar transaksi dengan filter opsional.
    - **GET** `/transactions/{id}`: Mendapatkan detail transaksi berdasarkan ID.
    - **PUT** `/transactions/{id}`: Memperbarui transaksi berdasarkan ID.
    - **DELETE** `/transactions/{id}`: Menghapus transaksi berdasarkan ID.

2. **Manajemen Kategori:**
    - **POST** `/categories`: Menambahkan kategori baru.
    - **GET** `/categories`: Mendapatkan daftar kategori.
    - **GET** `/categories/{id}`: Mendapatkan detail kategori berdasarkan ID.
    - **PUT** `/categories/{id}`: Memperbarui kategori berdasarkan ID.
    - **DELETE** `/categories/{id}`: Menghapus kategori berdasarkan ID.

3. **Autentikasi:**
    - **POST** `/auth/register`: Mendaftarkan pengguna baru.
    - **POST** `/auth/login`: Login dan mendapatkan token JWT.

### Contoh Permintaan dan Respons

- **Menambahkan Transaksi (POST `/transactions`)**

    Permintaan:
    ```json
    {
        "amount": 50000,
        "type": "income",
        "category_id": "60a7dff2b8c9b5bdf8e2e4d8",
        "description": "Gaji bulan April"
    }
    ```

    Respons:
    ```json
    {
        "id": "60b8d7c4b8c9b5bdf8e2e4e1",
        "amount": 50000,
        "type": "income",
        "category_id": "60a7dff2b8c9b5bdf8e2e4d8",
        "description": "Gaji bulan April",
        "created_at": "2024-04-01T12:34:56Z"
    }
    ```

- **Mengambil Daftar Transaksi (GET `/transactions`)**

    Permintaan:
    ```http
    GET /transactions?start_date=2024-01-01&end_date=2024-01-31 HTTP/1.1
    Host: localhost:8080
    Authorization: Bearer {token}
    ```

    Respons:
    ```json
    [
        {
            "id": "60b8d7c4b8c9b5bdf8e2e4e1",
            "amount": 50000,
            "type": "income",
            "category_id": "60a7dff2b8c9b5bdf8e2e4d8",
            "description": "Gaji bulan April",
            "created_at": "2024-04-01T12:34:56Z"
        },
        {
            "id": "60b8d7e2b8c9b5bdf8e2e4e2",
            "amount": -15000,
            "type": "expense",
            "category_id": "60a7dff2b8c9b5bdf8e2e4d9",
            "description": "Makan malam",
            "created_at": "2024-04-02T19:20:30Z"
        }
    ]
    ```

## Kontribusi 🤝

Kami menyambut kontribusi dari siapa saja. Jika Anda menemukan bug atau memiliki saran untuk fitur baru, silakan buat *issue* atau kirim *pull request*.

### Langkah Kontribusi

1. Fork repositori ini.
2. Buat *feature branch* (`git checkout -b feature/feature_name`).
3. Commit perubahan Anda (`git commit -m 'Add some feature'`).
4. Push ke branch (`git push origin feature/feature_name`).
5. Buat *pull request*.

## Lisensi 📜

Proyek ini dilisensikan di bawah lisensi MIT - lihat file [LICENSE](LICENSE) untuk detailnya.
