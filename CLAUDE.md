# YAVA — Claude Code Guide

**YAVA** (Your Amazing Various Aromas) — web app untuk menyimpan, mengelola, dan berbagi resep kopi.

---

## Tech Stack

| Layer | Tech |
|---|---|
| Frontend | Next.js (App Router) + TypeScript |
| Backend | Golang (REST API) |
| Database | PostgreSQL |
| Auth | Google OAuth 2.0 → JWT (HttpOnly cookie) |
| File Storage | Cloudflare R2 (S3-compatible) |
| Realtime | WebSocket / SSE |
| Frontend Deploy | Vercel |
| Backend Deploy | VPS + Docker |

---

## Monorepo Structure (Target)

```
yava/
├── apps/
│   ├── web/          # Next.js frontend
│   └── api/          # Golang backend
├── packages/
│   └── types/        # Shared TypeScript types (optional)
├── PRD.md
└── CLAUDE.md
```

---

## Design Reference

Prototype: `YAVA+Prototype.html` — tersimpan di claude.ai design share.
Link: https://claude.ai/design/p/cccb1e57-1e8f-4d2f-b476-0cfa68741af2?file=YAVA+Prototype.html&via=share

> Catatan: Design menggunakan warna dan komponen yang harus dikonsistenkan dengan Tailwind CSS + shadcn/ui.

---

## Task Breakdown

### PHASE 0 — Project Setup

- [ ] **P0-01** Init monorepo: buat folder `apps/web` dan `apps/api`
- [ ] **P0-02** Setup Next.js App Router di `apps/web` (TypeScript, Tailwind CSS, shadcn/ui)
- [ ] **P0-03** Setup Golang project di `apps/api` (Go modules, folder structure: `cmd/`, `internal/`, `pkg/`)
- [ ] **P0-04** Buat `docker-compose.yml` untuk PostgreSQL lokal
- [ ] **P0-05** Setup `.env.example` untuk web dan api
- [ ] **P0-06** Setup ESLint + Prettier untuk web
- [ ] **P0-07** Setup `golangci-lint` untuk api

---

### PHASE 1 — Auth + DB Schema + Recipe CRUD Dasar (Week 1–2)

#### Backend (Golang)

- [ ] **P1-01** Buat DB schema migration (tabel: `users`, `recipe_types`, `recipe_subtypes`, `recipes`, `recipe_sessions`, `recipe_notes`)
- [ ] **P1-02** Buat semua DB indexes (lihat PRD §12.3)
- [ ] **P1-03** Implementasi Google OAuth handler: `GET /auth/google` + `GET /auth/google/callback`
- [ ] **P1-04** Implementasi JWT issue + store ke HttpOnly cookie
- [ ] **P1-05** Implementasi `POST /auth/logout`
- [ ] **P1-06** Buat JWT middleware untuk semua protected routes
- [ ] **P1-07** Endpoint `GET /types` — list semua recipe types
- [ ] **P1-08** Endpoint `GET /types/:id/subtypes` — list subtypes
- [ ] **P1-09** Endpoint `POST /recipes` — create recipe (dengan sessions + notes)
- [ ] **P1-10** Endpoint `GET /recipes` — list recipes (query: `visibility`, `type_id`, `page`, `limit`)
- [ ] **P1-11** Endpoint `GET /recipes/:id` — detail recipe (sessions + notes ordered)
- [ ] **P1-12** Endpoint `PUT /recipes/:id` — update recipe (owner only)
- [ ] **P1-13** Seed data: recipe types (V60, Espresso, dll) + subtypes + default recipes

#### Frontend (Next.js)

- [ ] **P1-14** Setup layout utama: sidebar/navbar, auth guard
- [ ] **P1-15** Halaman login: tombol "Login with Google"
- [ ] **P1-16** Handle OAuth callback, simpan session state
- [ ] **P1-17** Dashboard: list resep milik user
- [ ] **P1-18** Halaman create recipe (multi-step form: info dasar → sessions → notes → visibility)
- [ ] **P1-19** Halaman detail recipe (tampilkan params, sessions, notes)
- [ ] **P1-20** Halaman edit recipe

---

### PHASE 2 — Timer System + Brewing Mode (Week 3–4)

- [ ] **P2-01** Komponen `BrewingTimer`: countdown timer menggunakan `performance.now()` (akurasi ±100ms)
- [ ] **P2-02** Logika auto-advance antar session (zero delay)
- [ ] **P2-03** Preparation countdown 3 detik sebelum session pertama ("Siapkan peralatan Anda...")
- [ ] **P2-04** Controls: Pause, Resume, Reset, Skip Session
- [ ] **P2-05** Progress bar per session
- [ ] **P2-06** Display `RecipeNote` di posisi order-nya (tanpa blok timer)
- [ ] **P2-07** Brewing Complete Screen: tampilkan total waktu, tombol "Ulangi" + "Kembali ke Resep"
- [ ] **P2-08** Full-screen brewing mode (mobile-friendly)
- [ ] **P2-09** Handle resep tanpa session (notes-only → checklist mode)

---

### PHASE 3 — Recipe Visibility + Explore + Duplicate (Week 5–6)

#### Backend

- [ ] **P3-01** Endpoint `PATCH /recipes/:id/archive`
- [ ] **P3-02** Endpoint `PATCH /recipes/:id/restore`
- [ ] **P3-03** Endpoint `POST /recipes/:id/duplicate`
- [ ] **P3-04** Filter list endpoint untuk public recipes (explore page)

#### Frontend

- [ ] **P3-05** Explore page: list public recipes + filter by type
- [ ] **P3-06** UI untuk archive/restore recipe (kebab menu di detail page)
- [ ] **P3-07** UI untuk duplicate recipe
- [ ] **P3-08** Badge "Archived" di dashboard untuk resep yang diarsip
- [ ] **P3-09** Badge "Default" untuk system recipes
- [ ] **P3-10** Visibility selector di form create/edit (Private / Public / Group)

---

### PHASE 4 — Group System (Week 7–9)

#### Backend

- [ ] **P4-01** DB migration: tabel `groups`, `group_members`, `group_recipes`
- [ ] **P4-02** Endpoint `POST /groups` — create group
- [ ] **P4-03** Endpoint `GET /groups/:id` — group detail
- [ ] **P4-04** Endpoint `PUT /groups/:id` — update group (admin only)
- [ ] **P4-05** Endpoint `DELETE /groups/:id` — disband group (founder only)
- [ ] **P4-06** Endpoint `GET /groups/:id/members` — list members
- [ ] **P4-07** Endpoint `POST /groups/:id/members` — join via invite code
- [ ] **P4-08** Endpoint `DELETE /groups/:id/members/:uid` — remove member (admin)
- [ ] **P4-09** Endpoint `PATCH /groups/:id/members/:uid/role` — ubah role (admin)
- [ ] **P4-10** Endpoint `GET /groups/:id/recipes` — list active group recipes
- [ ] **P4-11** Endpoint `GET /groups/:id/recipes/pending` — list pending (admin)
- [ ] **P4-12** Endpoint `POST /groups/:id/recipes` — submit recipe ke group
- [ ] **P4-13** Endpoint `PATCH /groups/:id/recipes/:rid/approve` — approve (admin)
- [ ] **P4-14** Endpoint `PATCH /groups/:id/recipes/:rid/reject` — reject + reason (admin)
- [ ] **P4-15** Endpoint `DELETE /groups/:id/recipes/:rid` — remove from group (admin)
- [ ] **P4-16** Generate unique `invite_code` saat group dibuat

#### Frontend

- [ ] **P4-17** Halaman create group
- [ ] **P4-18** Halaman group detail: tabs (Resep / Members / Settings)
- [ ] **P4-19** UI join group via invite link
- [ ] **P4-20** UI submit recipe ke group (dari detail page recipe)
- [ ] **P4-21** UI approve/reject pending recipes (admin view)
- [ ] **P4-22** UI manage members: list, remove, promote
- [ ] **P4-23** UI copy invite link

---

### PHASE 5 — Discussions + Notifications (Week 10–11)

#### Backend

- [ ] **P5-01** DB migration: tabel `discussions`, `notifications`
- [ ] **P5-02** Endpoint `GET /groups/:id/recipes/:rid/discussions`
- [ ] **P5-03** Endpoint `POST /groups/:id/recipes/:rid/discussions`
- [ ] **P5-04** Endpoint `PATCH /groups/:id/recipes/:rid/discussions/:did/pin` (admin)
- [ ] **P5-05** Endpoint `DELETE /groups/:id/recipes/:rid/discussions/:did`
- [ ] **P5-06** Endpoint `GET /notifications`
- [ ] **P5-07** Endpoint `PATCH /notifications/:id/read`
- [ ] **P5-08** Endpoint `PATCH /notifications/read-all`
- [ ] **P5-09** Trigger notifikasi: recipe approved/rejected, discussion reply
- [ ] **P5-10** SSE atau WebSocket endpoint untuk push notifikasi realtime

#### Frontend

- [ ] **P5-11** Komponen discussion thread (dengan nested replies)
- [ ] **P5-12** UI pin/unpin comment (admin)
- [ ] **P5-13** Notification bell di navbar: badge unread count
- [ ] **P5-14** Notification dropdown/panel: list notifikasi, mark as read

---

### PHASE 6 — QA + Deploy (Week 12)

- [ ] **P6-01** E2E test: auth flow
- [ ] **P6-02** E2E test: create + brew recipe
- [ ] **P6-03** E2E test: group submission + approval flow
- [ ] **P6-04** Performance test: API response time < 500ms
- [ ] **P6-05** Security audit: JWT validation, ownership checks, input sanitization
- [ ] **P6-06** Rate limiting: `/auth/*` 10 req/min, `POST /recipes` 30 req/min
- [ ] **P6-07** Setup Vercel deployment untuk web
- [ ] **P6-08** Setup Docker + VPS deployment untuk api
- [ ] **P6-09** Setup Cloudflare R2 bucket + upload profile photo
- [ ] **P6-10** Staging → Production sign-off

---

## API Conventions

- Base URL: `http://localhost:8080/v1` (dev), `https://api.yava.app/v1` (prod)
- Auth: `Authorization: Bearer <jwt>` (atau HttpOnly cookie)
- Response format selalu:
  ```json
  { "success": true, "data": {}, "meta": { "page": 1, "limit": 20, "total": 100 } }
  ```
- Error format:
  ```json
  { "success": false, "error": { "code": "RECIPE_NOT_FOUND", "message": "..." } }
  ```
- Semua list endpoint wajib pagination (`limit` default 20, max 100)

---

## Key Business Rules (Wajib Diingat)

1. **Timer**: auto-advance antar session = zero delay. Gunakan `performance.now()` bukan `Date.now()`.
2. **Default recipes**: `owner_id = NULL`, tidak bisa diedit/dihapus user biasa. Hanya bisa diduplikasi.
3. **Archive ≠ Delete**: resep yang diarsip tetap ada, cuma tersembunyi dari publik/grup.
4. **Group recipe flow**: submit → pending → approved/rejected. Owner tetap punya resep meski ditolak.
5. **RecipeSession.Order + RecipeNote.Order** berbagi sequence yang sama.
6. **Founder** adalah satu-satunya yang bisa membubarkan grup (`DELETE /groups/:id`).
7. **Visibility transition**: `private → public → archived` (tidak bisa balik ke private setelah public).

---

## Naming Conventions

- **Backend (Go)**: snake_case untuk DB, PascalCase untuk struct, camelCase untuk JSON response
- **Frontend (TS)**: camelCase untuk variabel/fungsi, PascalCase untuk komponen
- **API paths**: kebab-case (`/recipe-types` bukan `/recipeTypes`)
- **Branch**: `feat/P1-03-google-oauth`, `fix/P2-01-timer-drift`
- **Go files**: prefix per layer — `C_` Controller, `S_` Service, `Rp_` Repository, `M_` Model, `Rs_` Resource

---

## Backend Architecture (`apps/api`)

Framework: Fiber v2 + GORM + Goose (migrations). Semua operasi write **wajib** dalam transaction.

### Layer Structure

| Layer | Folder | Fungsi |
|---|---|---|
| Controller | `app/controllers/` | Parse request, panggil Service, return response |
| Service | `app/services/` | Logika bisnis & validasi |
| Repository | `app/repositories/` | Query builder ke DB |
| Model | `app/models/` | Struct tabel + implementasi `CoreModels` |
| Resource | `app/resources/` | Transform data ke format response |
| Middleware | `app/middlewares/` | JWT auth, rate limiter, locale, DB context |
| Routes | `routes/` | Definisi endpoint per domain |

### Request Flow

```
Request → Middleware → Controller → Service → Repository → Model
                                                         ← Repository
                                 ← Service (resource transform)
       ← Response JSON
```

### Transaction Rule

```go
if txErr := tx.Transaction(func(tx *gorm.DB) error {
    return tx.Create(&m).Error
}); txErr != nil {
    return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "gagal")
}
```

### Route Pattern

```
GET|POST|PUT|DELETE /api/admin/:locale/{resource}
```

`:locale` = `id` atau `en` untuk lokalisasi pesan. Route protected pakai `middlewares.CheckUserToken()`.

### API Response Format (framework)

```json
// List: { "data": [...], "meta": { "page": 1, "limit": 10, "total": 100, "total_pages": 10 } }
// Single: { "data": { "id": "1", ... } }
// Error: { "status": 422, "message": "...", "errors": { "field": ["msg"] } }
// Delete: { "data": { "message": "Deleted successfully" } }
```

### Referensi Implementasi

Lihat `Example` sebagai full-stack pattern: `routes/Example.go` → `C_Example.go` → `S_Example.go` → `Rp_Example.go` → `M_Example.go` → `Rs_Example.go`

---

## File Upload (Cloudflare R2)

- Profile photo: upload langsung dari browser ke presigned URL
- Step images (future): sama, presigned URL dari API
- Format yang didukung: JPG, PNG, WebP. Max 5MB per file.
