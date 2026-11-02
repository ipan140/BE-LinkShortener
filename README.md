# Ultra-Fast Redirector & Click Collector (Modular Monolith)

Backend performa ultra-tinggi untuk link shortener yang mengarahkan klik ke target URL dalam tempo **<10 ms** menggunakan **Redis In-Memory Cache**, sambil mengumpulkan data klik (Device, OS, Browser, IP, Negara, Kota, Referrer) secara non-blocking ke **Antrean Background (Redis Queue)** untuk diproses oleh Background Worker.

---

## 🏛️ Arsitektur Sistem

### 1. Pola Utama: Modular Monolith Architecture
* **Single Deployable Unit**: Seluruh sistem dikompilasi menjadi satu binary executable mandiri, mengeliminasi overhead jaringan internal dan kompleksitas orkestrasi microservices (RAM < 50 MB pada beban standar).
* **Domain Encapsulation**: Kode diorganisasikan per domain bisnis terpisah pada direktori `app/modules/` (`auth`, `link`, `redirect`, `analytics`). Tiap modul memiliki batas tanggung jawab (*bounded context*) yang terisolasi.

### 2. Pola Internal Modul: Clean Layered Architecture (N-Tier)
Setiap modul mengadopsi prinsip pemisahan tanggung jawab (*Separation of Concerns*) dengan struktur 5 lapisan:

```
[ HTTP Request dari Client / Browser ]
               │
               ▼
    ┌───────────────────────────┐
    │     route.go & handler    │  --> Presentation Layer (Routing Echo, Binding, Response)
    └──────────────┬────────────┘
                   │
                   ▼
    ┌───────────────────────────┐
    │        service.go         │  --> Business Logic Layer (Resolve URL, Validasi, Async Trigger)
    └──────────────┬────────────┘
                   │
                   ▼
    ┌───────────────────────────┐
    │       repository.go       │  --> Data Access Layer (GORM DB Queries, Pagination)
    └──────────────┬────────────┘
                   │
                   ▼
    ┌───────────────────────────┐
    │    model.go & redis.go    │  --> Entity Layer (GORM Struct) & Redis High-Speed Caching
    └───────────────────────────┘
```

* **route.go**: Mendaftarkan URL endpoint dan menghubungkan middleware ke route group Echo.
* **handler.go**: Memvalidasi request masuk, memanggil service, dan mengembalikan HTTP response / redirect.
* **service.go**: Menangani aturan bisnis murni tanpa keterikatan erat pada framework HTTP.
* **repository.go**: Mengabstraksi kueri basis data GORM, transaksi, dan agregasi data.
* **model.go**: Mendefinisikan struktur tabel relasional, foreign key, dan schema mapping.
* **redis.go**: Menyediakan strategi cache berkecepatan tinggi dan antrean event berbasis Redis.

---

## 🚀 Alur Kerja Ultra-Fast Redirect (<10 ms) & Click Collector

```
Browser Request (GET /:code)
       │
       ▼
[Echo HTTP Handler] ───────────────► Ekstrak IP, Referer, User-Agent, Lokasi
       │
       ▼
[Redis Cache Check] ──(Hit: <1ms)───► [302 Redirect dikirim ke User]
       │                                       │
       │ (Cache Miss)                          ▼
       ▼                             [Async Goroutine: Push Click Event ke Redis List]
[DB Fallback (GORM)]                                   │
       │                                               ▼
       └──► [Simpan ke Redis Cache (TTL 24h)]    [Background Worker: Batch Drain (100 item/1s)]
                                                       │
                                                       ▼
                                                 [Bulk Insert ke Database Clicks Table]
```

1. **Ultra-Fast Redirect Path**: Request `GET /:code` langsung membaca key `link:{code}` di Redis (waktu respon sub-milidetik). Browser menerima HTTP 302 Found seketika (<10 ms).
2. **Non-Blocking Click Collector**: Data klik (IP, browser, OS, device, referrer, geo) di-*serialize* dan di-*push* ke antrean Redis (`queue:clicks`) dalam goroutine terisolasi tanpa memblokir HTTP redirect.
3. **Batch Collector Worker**: Worker latar belakang mengambil event klik dalam *batch* (100 event atau per 1 detik) dan mengeksekusi *bulk insert* ke tabel relasional serta meng-update agregat klik secara efisien.

---

## 📂 Struktur Direktori

```text
├── app/
│   └── modules/
│       ├── auth/             # Modul Otentikasi & Akun Pengguna
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   └── service.go
│       ├── link/             # Modul Manajemen Short Link (CRUD & Cache Sync)
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   └── service.go
│       ├── redirect/         # Modul Ultra-Fast Redirector & Click Collector
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   └── service.go
│       └── analytics/        # Modul Background Worker & Statistik Klik
│           ├── handler.go
│           ├── model.go
│           ├── redis.go
│           ├── repository.go
│           ├── route.go
│           └── service.go
├── cmd/
│   └── api/
│       └── main.go           # Single Binary Entrypoint & DI Bootstrap
├── common/
│   ├── errors/               # Standardized Error Definitions
│   ├── middleware/           # JWT, CORS, High-Performance Logger
│   └── utils/                # Base62 Generator, Real IP, UserAgent Parser
├── config/
│   ├── config.go             # App, DB, Redis, Worker Configurations
│   ├── database.go           # GORM Postgres/SQLite Connection Pooling
│   ├── env.go                # Env Loader & Helpers
│   └── redis.go              # Go-Redis Client & Connection Pool
├── docs/
│   ├── docs.go               # Swagger UI Embedded Handler
│   └── openapi.json          # OpenAPI 3.0 Specification
├── .dockerignore
├── .editorconfig
├── .env.example
├── .gitattributes
├── .gitignore
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

---

## ⚙️ Persyaratan Sistem

* **Go**: Versi 1.24+ (atau versi Go terbaru yang terpasang)
* **Redis**: Versi 6.2+ (untuk *in-memory cache* dan *queue*)
* **Database**: SQLite (default lokal) atau PostgreSQL (production)

---

## 🛠️ Panduan Instalasi & Menjalankan

### 1. Salin Environment
```bash
cp .env.example .env
```

### 2. Jalankan Server Secara Langsung
```bash
go run ./cmd/api
```

### 3. Build Single Executable Binary
```bash
go build -ldflags="-s -w" -o redirector.exe ./cmd/api
```

### 4. Menggunakan Docker
```bash
docker build -t redirector:latest .
docker run -p 8080:8080 --env-file .env redirector:latest
```

---

## 📡 Daftar Endpoint API

### 1. Ultra-Fast Redirect (Publik)
* `GET /:code` -> Redirect langsung ke tujuan target (HTTP 302 Found)
* `GET /r/:code` -> Alternatif endpoint redirect langsung

### 2. Autentikasi (`/api/v1/auth`)
* `POST /api/v1/auth/register` -> Mendaftar pengguna baru
* `POST /api/v1/auth/login` -> Login & peroleh JWT token
* `GET /api/v1/auth/me` -> Info profil (Wajib Bearer Token)

### 3. Manajemen Short Link (`/api/v1/links`)
* `POST /api/v1/links` -> Buat tautan pendek baru (Auto/Custom Alias, TTL)
* `GET /api/v1/links` -> Daftar tautan pengguna (Paginated)
* `GET /api/v1/links/:id` -> Detail tautan
* `PUT /api/v1/links/:id` -> Update tautan & sinkronisasi cache
* `DELETE /api/v1/links/:id` -> Hapus tautan & hapus cache

### 4. Analitik & Koleksi Klik (`/api/v1/analytics`)
* `GET /api/v1/analytics/:code` -> Ringkasan klik, browser, OS, perangkat, negara, referrer, dan grafik harian
