# CLAUDE.md — Go Backend Framework

Panduan arsitektur dan konvensi untuk framework backend ini.

---

## Arsitektur Layer

| Layer | Folder | Fungsi |
|---|---|---|
| Controller | `app/controllers/` | Terima HTTP request, panggil Service, kembalikan response |
| Service | `app/services/` | Logika bisnis, validasi input, orkestrasi proses |
| Repository | `app/repositories/` | Query builder & custom query ke database |
| Model | `app/models/` | Struct tabel DB + implementasi `CoreModels` interface |
| Resource | `app/resources/` | Transform data ke format response API |
| Middleware | `app/middlewares/` | JWT auth, rate limiter, locale, DB context |
| Routes | `routes/` | Definisi endpoint dan pengelompokan route |
| Helpers | `helpers/` | Utility: converter, validator, JWT, bcrypt, dsb |
| Lang | `lang/` | Pesan lokalisasi (error message dalam id/en/dll) |
| Exceptions | `exceptions/` | Format standar error response |
| Cmd | `cmd/migrate/` | Perintah migrasi database (goose) |
| Database | `database/` | Koneksi DB, seed, dan setup |
| Public | `public/` | Static files (img, template) |
| Docs | `docs/` | Dokumentasi API (Swagger, dsb) |

---

## Naming Convention

| Tipe | Prefix | Contoh |
|---|---|---|
| Controller | `C_` | `C_Example.go` |
| Service | `S_` | `S_Example.go` |
| Repository | `Rp_` | `Rp_Example.go` |
| Model | `M_` | `M_Example.go` |
| Resource | `Rs_` | `Rs_Example.go` |

---

## Route Pattern

```
GET|POST|PUT|DELETE /api/admin/:locale/{resource}
```

- `:locale` diisi `id` atau `en` — digunakan untuk lokalisasi pesan
- Route dikelompokkan per domain dalam `routes/`
- Route protected: tambahkan `middlewares.CheckUserToken()` pada group

---

## Alur Request

```
Request
  → Middleware (JWT auth, locale, DB context)
    → Controller (parse body/params)
      → Service (validasi, logika bisnis)
        → Repository (query DB)
          → Model (struct + scopes)
        ← Repository
      ← Service
    ← Controller (resource transform)
  → Response JSON
```

---

## Transaction Rule

Semua operasi write ke DB (`Create`, `Update`, `Delete`) **wajib** dalam `tx.Transaction()`.
Framework sudah memasang callback GORM yang memblokir write di luar transaction.

```go
if txErr := tx.Transaction(func(tx *gorm.DB) error {
    return tx.Create(&m).Error
}); txErr != nil {
    return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "gagal membuat data")
}
```

---

## Contoh Lengkap

Lihat implementasi `Example` sebagai referensi full-stack pattern:

- Route: `routes/Example.go`
- Controller: `app/controllers/C_Example.go`
- Service: `app/services/S_Example.go`
- Repository: `app/repositories/Rp_Example.go`
- Model: `app/models/M_Example.go`
- Resource: `app/resources/Rs_Example.go`

---

## Environment Variables

Buat file `.env` dari template berikut:

```env
PORT=:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=secret
DB_DATABASE=mydb
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:5173
```

---

## Response Format

**List (paginated):**
```json
{
  "data": [...],
  "meta": { "page": 1, "limit": 10, "total": 100, "total_pages": 10 }
}
```

**Single:**
```json
{ "data": { "id": "1", "name": "..." } }
```

**Error:**
```json
{ "status": 422, "message": "...", "errors": { "field": ["msg"] } }
```

**Delete:**
```json
{ "data": { "message": "Deleted successfully" } }
```
