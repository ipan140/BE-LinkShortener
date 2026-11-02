# Ultra-Fast Redirector & Click Collector (Modular Monolith)

Backend performa ultra-tinggi untuk link shortener yang mengarahkan klik ke target URL dalam tempo **<10 ms** menggunakan **Redis In-Memory Cache**, sambil mengumpulkan data klik (Device, OS, Browser, IP, Negara, Kota, Referrer, UTM Tags) secara non-blocking ke **Antrean Background (Redis Queue)** untuk diproses oleh Background Worker ke **PostgreSQL**.

---

## 🌟 Fitur Unggulan

1. **Ultra-Fast Redirector (<10 ms Target)**:
   * Query in-memory cache langsung dari Redis (`link:{code}`) dengan latensi sub-milidetik.
   * Auto-fallback ke database PostgreSQL GORM jika terjadi cache miss, sekaligus melakukan *cache warming* (TTL 24 jam).
2. **📱 Smart Device-Based Routing (Deep-Linking iOS / Android)**:
   * Pengunjung dari perangkat **iOS** (iPhone/iPad) otomatis diarahkan ke `ios_url` (App Store).
   * Pengunjung dari perangkat **Android** otomatis diarahkan ke `android_url` (Google Play Store).
   * Pengunjung Desktop/Web diarahkan ke `original_url`.
3. **🏷️ UTM & Query Parameter Forwarding**:
   * Secara otomatis meneruskan dan menggabungkan query parameter (seperti `utm_source`, `utm_medium`, `utm_campaign`, `fbclid`, dll.) dari shortlink ke URL tujuan akhir tanpa merusak parameter yang sudah ada di URL asli.
4. **📊 Export Data Analitik ke CSV (`GET /api/v1/analytics/:code/export.csv`)**:
   * Unduh seluruh riwayat telemetri klik secara langsung dalam format spreadsheet `.csv` (Timestamp UTC, IP, Negara, Kota, Device, Browser, OS, Referrer).
5. **🎨 Responsive HTML Landing Pages (404, Expired, Quota Reached, Deactivated)**:
   * Jika link tidak ditemukan atau sudah kedaluwarsa, browser pengunjung disajikan halaman HTML modern dark-theme interaktif (bukan JSON mentah). Klien API tetap mendapatkan JSON standar.
6. **🔒 Password-Protected Links**:
   * Mengamankan link dengan kata sandi (menggunakan enkripsi kuat bcrypt) dengan formulir PIN interaktif di browser.
7. **⏱️ Max Click Limit (Self-Destruct)**:
   * Batasi jumlah total klik yang diizinkan (misal kuota 100 klik). Setelah limit tercapai, link otomatis kedaluwarsa.
8. **🛡️ Anti-Abuse Redis Sliding-Window Rate Limiter**:
   * Melindungi endpoint redirect dan REST API dari serangan DDoS, click fraud, dan bot flooding menggunakan algoritma Redis sliding-window ZSet.
9. **📷 Instant QR Code Generator**:
   * Endpoint publik `GET /qr/:code?size=256` untuk menghasilkan gambar QR Code PNG beresolusi tinggi langsung dari memory.

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
│       ├── link/              # Modul Manajemen Link, Device Routing & QR Generator
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   └── service.go
│       ├── redirect/          # Ultra-Fast Redirector, Smart Routing & Landing Pages
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── redis.go
│       │   ├── repository.go
│       │   ├── route.go
│       │   ├── service.go
│       │   └── service_test.go
│       └── analytics/         # Worker Telemetri, Agregat & Export CSV
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
│   ├── database.go            # GORM PostgreSQL Connection Pooling
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

### 1. Menjalankan Server
```bash
go run ./cmd/api
```

### 2. Menjalankan Unit Test
```bash
go test -v ./...
```

### 3. Build Single Executable Binary
```bash
go build -ldflags="-s -w" -o redirector.exe ./cmd/api
```

### 4. Swagger UI Documentation
Akses dokumentasi interaktif di:
👉 **`http://localhost:8080/swagger/index.html`**
