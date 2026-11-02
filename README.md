# Ultra-Fast Redirector & Click Collector (Modular Monolith)

Backend performa ultra-tinggi untuk link shortener yang mengarahkan klik ke target URL dalam tempo **<10 ms** menggunakan **Redis In-Memory Cache**, sambil mengumpulkan data klik (Device, OS, Browser, IP, Negara, Kota, Referrer, UTM Tags) secara non-blocking ke **Antrean Background (Redis Queue)** untuk diproses oleh Background Worker ke **PostgreSQL**.

---

## 🌟 Fitur Lengkap Sistem

1. **Ultra-Fast Redirector (<10 ms Target)**:
   * Query in-memory cache langsung dari Redis (`link:{code}`) dengan latensi sub-milidetik.
   * Auto-fallback ke database PostgreSQL GORM jika terjadi cache miss, sekaligus melakukan *cache warming* (TTL 24 jam).
2. **🔔 Webhook Notifications (Discord, Telegram, Custom HTTP)**:
   * Notifikasi otomatis ke URL Webhook saat shortlink mencapai milestone klik (10, 50, 100, 500, 1000, 5000, dst.) atau saat batas kuota `max_clicks` tercapai.
   * Format otomatis **Discord Rich Embed** dengan warna status khusus saat URL webhook berasal dari `discord.com/api/webhooks`.
3. **📱 Smart Device-Based Routing (Deep-Linking iOS / Android)**:
   * Pengunjung dari perangkat **iOS** (iPhone/iPad) otomatis diarahkan ke `ios_url` (App Store).
   * Pengunjung dari perangkat **Android** otomatis diarahkan ke `android_url` (Google Play Store).
   * Pengunjung Desktop/Web diarahkan ke `original_url`.
4. **🌐 Social Media Rich Previews (OpenGraph & Twitter Cards)**:
   * Deteksi bot media sosial (WhatsApp, Telegram, Twitter/X, Discord, Slack, Facebook).
   * Menampilkan kartu pratinjau judul, deskripsi, dan banner gambar (`image_url`) sebelum diklik.
5. **🔑 Permanent API Key Authentication (`X-API-Key`)**:
   * Pengguna mendapatkan API Key permanen (`ak_live_...`) saat mendaftar. Mendukung integrasi tanpa kedaluwarsa dengan script otomatisasi via header `X-API-Key`.
6. **🏷️ Kategori & Tag Kampanye (`#tag`)**:
   * Mengelompokkan tautan berdasarkan tag (`tags: "marketing,promo"`), serta memfilter daftar tautan melalui `GET /api/v1/links?tag=promo`.
7. **📊 Export Data Analitik ke CSV (`GET /api/v1/analytics/:code/export.csv`)**:
   * Unduh seluruh riwayat telemetri klik secara langsung dalam format spreadsheet `.csv`.
8. **🎨 Responsive HTML Landing Pages (404, Expired, Quota Reached, Deactivated)**:
   * Kunjungan browser ke link yang kedaluwarsa/tidak ada akan disajikan halaman HTML modern bertema gelap.
9. **🔒 Password-Protected Links**:
   * Tautan dengan kata sandi (dienkripsi bcrypt) dan formulir PIN modern responsif.
10. **⏱️ Max Click Limit (Self-Destruct)**:
    * Pembatasan kuota klik total yang otomatis menonaktifkan link saat kuota tercapai.
11. **🛡️ Anti-Abuse Redis Sliding-Window Rate Limiter**:
    * Perlindungan terhadap serangan DDoS dan spam bot klik menggunakan Redis Sorted Set (`ZSet`).
12. **📷 Instant QR Code Generator**:
    * Endpoint publik `GET /qr/:code?size=256` untuk menghasilkan gambar QR Code PNG instan.

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
