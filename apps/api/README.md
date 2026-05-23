# Go Framework

Core backend framework berbasis [Go Fiber](https://gofiber.io/) + [GORM](https://gorm.io/), dirancang sebagai template siap pakai untuk semua project backend.

## Stack

- **Runtime**: Go 1.25+
- **Router**: Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL (branch `main`) — MySQL dan SQLite tersedia di branch terpisah
- **Migration**: [Goose](https://github.com/pressly/goose)
- **Auth**: JWT via httpOnly cookie / Authorization header

## Arsitektur

```
app/
├── controllers/   # HTTP handler — terima request, kembalikan response
├── services/      # Logika bisnis & validasi
├── repositories/  # Query builder & custom query
├── models/        # Struct tabel + GORM scopes
├── resources/     # Transform data ke format response
└── middlewares/   # JWT auth, rate limiter, locale
routes/            # Definisi endpoint per domain
helpers/           # Utility: converter, JWT, bcrypt, dsb
exceptions/        # Format error response standar
lang/              # Lokalisasi pesan (id, en)
database/
└── migrations/    # SQL migration (goose)
cmd/migrate/       # CLI migration runner
```

Lihat [CLAUDE.md](./CLAUDE.md) untuk panduan lengkap konvensi dan arsitektur.

## Setup

```bash
cp .env.example .env
# Edit .env sesuai konfigurasi DB

go mod download
go run cmd/migrate/main.go up
go run main.go
```

## Migration

```bash
go run cmd/migrate/main.go up          # jalankan semua migration
go run cmd/migrate/main.go down        # rollback 1 step
go run cmd/migrate/main.go status      # cek status migration
go run cmd/migrate/main.go create nama # buat migration baru
```

## Contoh Endpoint

```
GET    /api/admin/:locale/example        # list data
GET    /api/admin/:locale/example/:id    # detail
POST   /api/admin/:locale/example        # create
PUT    /api/admin/:locale/example/:id    # update
DELETE /api/admin/:locale/example/:id    # delete
```

Lihat implementasi di `routes/Example.go`, `app/controllers/C_Example.go`, dst sebagai referensi pattern.

## Branches

| Branch | Database |
|---|---|
| `main` | PostgreSQL |
| `mysql` | MySQL / MariaDB |
| `sqlite` | SQLite |
