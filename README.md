# Ultra-Fast Redirector & Click Collector (Modular Monolith)

Backend performa ultra-tinggi untuk link shortener yang mengarahkan klik ke target URL dalam tempo **<10 ms** menggunakan **Redis In-Memory Cache**, sambil mengumpulkan data klik (Device, OS, Browser, IP, Negara, Kota, Referrer, UTM Tags) secara non-blocking ke **Antrean Background (Redis Queue)** untuk diproses oleh Background Worker.

---

## 🌟 Fitur Unggulan

1. **Ultra-Fast Redirector (<10 ms Target)**:
   * Query in-memory cache langsung dari Redis (`link:{code}`) dengan latensi sub-milidetik.
   * Auto-fallback ke database GORM jika terjadi cache miss, sekaligus melakukan *cache warming* (TTL 24 jam).
2. **Asynchronous Click Collector**:
   * Pengumpulan metadata pengunjung (IP, negara via Cloudflare/Geo header, kota, browser, OS, device, referrer, query string) secara non-blocking melalui Redis list queue.
   * Background batch collector worker menguras antrean (100 item/1s) dan melakukan bulk insert ke basis data.
3. **QR Code Generator**:
   * Endpoint `GET /api/v1/links/qr/:code?size=256` untuk menghasilkan gambar QR Code PNG beresolusi tinggi langsung dari binary.
4. **UTM Parameter Forwarder**:
   * Secara otomatis meneruskan dan menggabungkan query parameter (seperti `utm_source`, `utm_medium`, `utm_campaign`, `fbclid`, dll.) dari shortlink ke URL tujuan akhir tanpa merusak parameter yang sudah ada di URL asli.
5. **Password-Protected Links**:
   * Mengamankan link dengan kata sandi (menggunakan enkripsi kuat bcrypt). Pengunjung via browser akan disuguhkan halaman PIN/Password interaktif modern, atau via query parameter `?p=password`.
6. **Max Click Limit (Self-Destruct)**:
   * Batasi jumlah total klik yang diizinkan (misal kuota 100 klik). Setelah limit tercapai, link otomatis kedaluwarsa.
7. **Anti-Abuse Redis Sliding-Window Rate Limiter**:
   * Melindungi endpoint redirect dan REST API dari serangan DDoS, click fraud, dan bot flooding menggunakan algoritma Redis sliding-window ZSet.

---

## 🏛️ Arsitektur Sistem

### 1. Pola Utama: Modular Monolith Architecture
* **Single Deployable Unit**: Seluruh sistem dikompilasi menjadi satu binary executable mandiri (`redirector.exe`), mengeliminasi overhead jaringan internal dan kompleksitas orkestrasi microservices (RAM < 50 MB pada beban standar).
* **Domain Encapsulation**: Kode diorganisasikan per domain bisnis terpisah pada direktori `app/modules/` (`auth`, `link`, `redirect`, `analytics`). Tiap modul memiliki batas tanggung jawab (*bounded context*) yang terisolasi.

### 2. Pola Internal Modul: Clean Layered Architecture (N-Tier)
Setiap modul mengadopsi prinsip pemisahan tanggung jawab (*Separation of Concerns*) dengan 5 lapisan:

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

---

## 📂 Struktur Direktori

```text
├── app/
│   └── modules/
│       ├── auth/              # Modul Otentikasi & Akun Pengguna (JWT)
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   └── service.go
│       ├── link/              # Modul Manajemen Link & QR Generator
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   └── service.go
│       ├── redirect/          # Ultra-Fast Redirector, UTM Forwarder & Protection
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   ├── service.go
│       │   └── service_test.go
│       └── analytics/         # Background Worker & Statistik Klik
│           ├── handler.go
│           ├── model.go
│           ├── redis.go
│           ├── repository.go
│           ├── route.go
│           └── service.go
├── cmd/
│   └── api/
│       └── main.go            # Monolith Entrypoint & DI Wiring
├── common/
│   ├── errors/                # Standardized HTTP Errors
│   ├── middleware/            # JWT, CORS, Logger, Redis Sliding-Window RateLimiter
│   └── utils/                 # Base62 Generator, IP Resolver, UserAgent Parser
├── config/
│   ├── config.go              # App, DB, Redis, Worker Configurations
│   ├── database.go            # GORM Postgres/SQLite Connection Pooling
│   ├── env.go                 # Env Loader & Helpers
│   └── redis.go               # Go-Redis Client & Connection Pool
├── docs/
│   ├── docs.go                # Embedded Swagger UI
│   └── openapi.json           # OpenAPI 3.0 Specification
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

## 🛠️ Panduan Menjalankan

### 1. Konfigurasi Lingkungan
Salin file template `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

### 2. Menjalankan Server
```bash
go run ./cmd/api
```

### 3. Menjalankan Unit Test
```bash
go test -v ./...
```

### 4. Build Single Executable Binary
```bash
go build -ldflags="-s -w" -o redirector.exe ./cmd/api
```

### 5. Swagger UI
Akses dokumentasi interaktif di: `http://localhost:8080/swagger/index.html`

---

## 📡 Panduan Penggunaan Fitur Baru

### A. Membuat Link dengan Proteksi Password & Batas Klik
`POST /api/v1/links`
```json
{
  "original_url": "https://company.internal/secret-report",
  "custom_code": "laporan-q3",
  "password": "PIN_RAHASIA_123",
  "max_clicks": 50,
  "title": "Laporan Rahasia Q3"
}
```

### B. Mengunduh QR Code Shortlink
`GET /api/v1/links/qr/laporan-q3?size=300`
* Mengembalikan file binary `image/png` yang siap dipajang di poster atau media sosial.

### C. Redirect dengan UTM Parameter Forwarding
Kunjungan ke:
`http://localhost:8080/laporan-q3?utm_source=telegram&utm_campaign=launch`
Akan secara otomatis diarahkan ke:
`https://company.internal/secret-report?utm_source=telegram&utm_campaign=launch`
Semua parameter UTM tetap utuh dan ikut tercatat di tabel log telemetri analitik!
