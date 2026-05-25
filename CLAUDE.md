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

> **Strategi**: FE dikerjakan lebih dulu dengan mock data untuk memvalidasi flow dan design. BE menyusul setelah FE per phase selesai.

---

## FRONTEND (Next.js — `apps/web`)

### PHASE 0 — Setup FE

- [x] **P0-01** Init monorepo: folder `apps/web` dan `apps/api`
- [x] **P0-02** Setup Next.js App Router (TypeScript, Tailwind CSS, shadcn/ui)
- [x] **P0-05** Setup `.env.example` untuk web
- [x] **P0-06** Setup ESLint + Prettier

### PHASE 1 — Auth + Recipe CRUD UI

- [x] **P1-14** Setup layout utama: topnav, route group `(app)`
- [x] **P1-15** Halaman login: tombol "Login with Google"
- [ ] **P1-16** Handle OAuth callback, simpan session state (tunggu BE P1-03)
- [x] **P1-17** Dashboard: list resep milik user (mock)
- [x] **P1-18** Halaman create recipe (multi-step form: info dasar → alur → visibilitas)
- [x] **P1-19** Halaman detail recipe (params, sessions, notes — mock)
- [x] **P1-20** Halaman edit recipe

### PHASE 2 — Timer + Brewing Mode UI

- [x] **P2-01** `BrewingTimer`: countdown + `performance.now()`
- [x] **P2-02** Auto-advance antar session (zero delay)
- [x] **P2-03** Preparation countdown 3 detik
- [x] **P2-04** Controls: Pause, Resume, Reset, Skip
- [x] **P2-05** Progress bar per session + total
- [x] **P2-06** Display `RecipeNote` di timeline (tanpa blok timer)
- [x] **P2-07** Brewing Complete Screen (total waktu, Ulangi, Kembali)
- [x] **P2-08** Full-screen brewing mode (`position: fixed`)
- [x] **P2-09** Handle resep tanpa session (notes-only → checklist mode)

### PHASE 3 — Explore + Archive/Duplicate UI

- [x] **P3-05** Explore page: list public recipes + search + filter by type
- [x] **P3-06** UI archive/restore recipe (kebab dropdown di detail page)
- [x] **P3-07** UI duplicate recipe (konfirmasi + toast)
- [x] **P3-08** Badge "ARSIP" + section terpisah di dashboard
- [x] **P3-09** Badge "DEFAULT" untuk system recipes (via `isDefault`)
- [x] **P3-10** Visibility selector di form create (Private / Public / Group)

### PHASE 4 — Group System UI

- [x] **P4-17** Halaman create group
- [x] **P4-18** Halaman group detail: tabs (Resep / Members / Settings)
- [x] **P4-19** UI join group via invite link
- [x] **P4-20** UI submit recipe ke group (dari detail page)
- [x] **P4-21** UI approve/reject pending recipes (admin view)
- [x] **P4-22** UI manage members: list, remove, promote
- [x] **P4-23** UI copy invite link

### PHASE 5 — Discussions + Notifikasi UI

- [x] **P5-11** Komponen discussion thread (nested replies)
- [x] **P5-12** UI pin/unpin comment (admin)
- [x] **P5-13** Notification bell: badge unread count
- [x] **P5-14** Notification dropdown: list + mark as read

### PHASE 6 — QA + Deploy FE

- [ ] **P6-02** E2E test: create + brew recipe flow
- [ ] **P6-07** Setup Vercel deployment untuk web

---

## BACKEND (Golang — `apps/api`)

### PHASE 0 — Setup BE

- [x] **P0-01** Init `apps/api` (Go modules, struktur folder)
- [x] **P0-03** Folder structure: `app/controllers/`, `app/services/`, `app/repositories/`, dll
- [x] **P0-04** `docker-compose.yml` untuk PostgreSQL lokal
- [x] **P0-05** Setup `.env.example` untuk api
- [ ] **P0-07** Setup `golangci-lint`

### PHASE 1 — Auth + Recipe CRUD API

- [x] **P1-01** DB schema migration: `yv_user`, `yv_cd_recipe_type`, `yv_cd_recipe_subtype`, `yv_recipe`, `yv_recipe_session`, `yv_recipe_note`
- [x] **P1-02** DB indexes (lihat PRD §12.3)
- [x] **P1-03** Google OAuth handler: `GET /v1/auth/google` + `GET /v1/auth/google/callback`
- [x] **P1-04** JWT issue + store ke HttpOnly cookie (`yv_token`)
- [x] **P1-05** `POST /v1/auth/logout`
- [x] **P1-06** JWT middleware — `YvAuth()` di `middlewares/YvAuth.go`
- [x] **P1-07** `GET /v1/types` — list recipe types
- [x] **P1-08** `GET /v1/types/:id/subtypes`
- [x] **P1-09** `POST /v1/recipes` — create (sessions + notes, dalam transaction)
- [x] **P1-10** `GET /v1/recipes` — list (filter: `visibility`, `type_id`, `mine`, pagination)
- [x] **P1-11** `GET /v1/recipes/:id` — detail (sessions + notes ordered)
- [x] **P1-12** `PUT /v1/recipes/:id` — update (owner only)
- [x] **P1-13** Seed: recipe types + subtypes + default recipes (migrations 05–07)

### PHASE 3 — Visibility + Archive API

- [ ] **P3-01** `PATCH /recipes/:id/archive`
- [ ] **P3-02** `PATCH /recipes/:id/restore`
- [ ] **P3-03** `POST /recipes/:id/duplicate`
- [ ] **P3-04** Filter public recipes untuk explore endpoint

### PHASE 4 — Group System API

- [ ] **P4-01** Migration: `groups`, `group_members`, `group_recipes`
- [ ] **P4-02** `POST /groups`
- [ ] **P4-03** `GET /groups/:id`
- [ ] **P4-04** `PUT /groups/:id` (admin)
- [ ] **P4-05** `DELETE /groups/:id` (founder)
- [ ] **P4-06** `GET /groups/:id/members`
- [ ] **P4-07** `POST /groups/:id/members` (join via invite code)
- [ ] **P4-08** `DELETE /groups/:id/members/:uid` (admin)
- [ ] **P4-09** `PATCH /groups/:id/members/:uid/role` (admin)
- [ ] **P4-10** `GET /groups/:id/recipes`
- [ ] **P4-11** `GET /groups/:id/recipes/pending` (admin)
- [ ] **P4-12** `POST /groups/:id/recipes` — submit
- [ ] **P4-13** `PATCH /groups/:id/recipes/:rid/approve`
- [ ] **P4-14** `PATCH /groups/:id/recipes/:rid/reject`
- [ ] **P4-15** `DELETE /groups/:id/recipes/:rid`
- [ ] **P4-16** Generate unique `invite_code`

### PHASE 5 — Discussions + Notifications API

- [ ] **P5-01** Migration: `discussions`, `notifications`
- [ ] **P5-02** `GET /groups/:id/recipes/:rid/discussions`
- [ ] **P5-03** `POST /groups/:id/recipes/:rid/discussions`
- [ ] **P5-04** `PATCH …/discussions/:did/pin` (admin)
- [ ] **P5-05** `DELETE …/discussions/:did`
- [ ] **P5-06** `GET /notifications`
- [ ] **P5-07** `PATCH /notifications/:id/read`
- [ ] **P5-08** `PATCH /notifications/read-all`
- [ ] **P5-09** Trigger notifikasi (approved/rejected/reply)
- [ ] **P5-10** SSE atau WebSocket untuk push notifikasi

### PHASE 6 — QA + Deploy BE

- [ ] **P6-01** E2E test: auth flow
- [ ] **P6-03** E2E test: group submission + approval flow
- [ ] **P6-04** Performance test: API response < 500ms
- [ ] **P6-05** Security audit: JWT, ownership checks, sanitasi input
- [ ] **P6-06** Rate limiting: `/auth/*` 10 req/min, `POST /recipes` 30 req/min
- [ ] **P6-08** Docker + VPS deployment
- [ ] **P6-09** Cloudflare R2 bucket + upload profile photo
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
- **DB tables**:
  - `yv_cd_{{name}}` → master data / jarang berubah (contoh: `yv_cd_client`, `yv_cd_setting`, `yv_cd_recipe_type`)
  - `yv_{{name}}` → data transaksional / sering diolah (contoh: `yv_recipe`, `yv_group`, `yv_notification`)

---

## Backend Architecture (`apps/api`)

Framework: Fiber v2 + GORM + Goose (migrations). Semua operasi write **wajib** dalam transaction.

### Layer Structure

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
| Database | `database/` | Koneksi DB global (`database.DB`) |
| Public | `public/` | Static files (img, template) |
| Docs | `docs/` | Dokumentasi API (Swagger, dsb) |

### Naming Convention (file prefix)

| Tipe | Prefix | Contoh |
|---|---|---|
| Controller | `C_` | `C_Example.go` |
| Service | `S_` | `S_Example.go` |
| Repository | `Rp_` | `Rp_Example.go` |
| Model | `M_` | `M_Example.go` |
| Resource | `Rs_` | `Rs_Example.go` |

### Request Flow

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

### Transaction Rule

Semua operasi write (`Create`, `Update`, `Delete`) **wajib** dalam `tx.Transaction()`. Framework sudah memasang callback GORM yang memblokir write di luar transaction.

```go
if txErr := tx.Transaction(func(tx *gorm.DB) error {
    return tx.Create(&m).Error
}); txErr != nil {
    return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "gagal membuat data")
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

### Template Wajib — Contoh Lengkap per Layer

Ketika membuat endpoint baru, **wajib** mengikuti pola ini. Ganti `Example`/`ExampleModel`/`tr_examples` dengan nama resource baru.

**routes/Example.go**
```go
func ExampleRoutes(r fiber.Router) {
    example := r.Group("example")
    example.Get("/", controllers.ExampleList)
    example.Get("/:id", controllers.ExampleShow)
    example.Post("/", controllers.ExampleCreate)
    example.Put("/:id", controllers.ExampleUpdate)
    example.Delete("/:id", controllers.ExampleDelete)
}
```

**app/controllers/C_Example.go** — parse → service → resource, tidak ada logika bisnis
```go
func ExampleList(c *fiber.Ctx) error {
    result, err := services.ExampleList(getDB(c), getBodyData(c), c, getLocale(c))
    if err != nil {
        return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
    }
    return c.JSON(resources.ExampleResource(c, result))
}
```

**app/services/S_Example.go** — validasi + logika bisnis + transaction
```go
func ExampleCreate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
    if _, err := models.Validate(c, data, new(ValidateExampleCreate), locale); err != nil {
        return nil, err
    }
    m := models.ExampleModel{Name: helpers.Conv(data["name"]).String()}
    if txErr := tx.Transaction(func(tx *gorm.DB) error {
        return tx.Create(&m).Error
    }); txErr != nil {
        return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "failed to create "+m.ModulName())
    }
    return repositories.ExampleSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
        return db.Where(models.ExampleModel{Id: m.Id})
    })
}
```

**app/repositories/Rp_Example.go** — gunakan helper `GetIndexData` / `GetSingleData` / `GetMultipleData`
```go
func ExampleIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (models.IndexData, any) {
    return models.GetIndexData(tx, data, c, locale, models.ExampleModel{})
}
func ExampleSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(db *gorm.DB) *gorm.DB) (map[string]any, any) {
    return models.GetSingleData(tx, data, c, locale, where, models.ExampleModel{})
}
```

**app/models/M_Example.go** — implementasi `CoreModels` (GetSelect, Searchable, Sortable, Join, Option)
```go
type ExampleModel struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}
func (ExampleModel) TableName() string { return "tr_examples" }
func (ExampleModel) ModulName() string { return "Example" }
func (s ExampleModel) ScopesGetSelect(data map[string]any) map[string]string {
    return map[string]string{"id": "tr_examples.id", "name": "tr_examples.name"}
}
func (s ExampleModel) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
    return map[string]SearchableFields{
        "id":   {Operators: []string{"=", "!="}},
        "name": {Operators: []string{"=", "like"}},
    }
}
func (s ExampleModel) ScopesSortbleFields(data map[string]any) map[string]bool {
    return map[string]bool{"id": true, "name": true}
}
func (s ExampleModel) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
    return func(tx *gorm.DB) *gorm.DB { return tx }
}
func (s ExampleModel) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
    return func(tx *gorm.DB) *gorm.DB { return tx }
}
func (s *ExampleModel) BeforeCreate(tx *gorm.DB) error { return AutoFillCreate(s, tx) }
func (s *ExampleModel) BeforeUpdate(tx *gorm.DB) error { return AutoFillUpdate(s, tx) }
func (s *ExampleModel) BeforeDelete(tx *gorm.DB) error { return AutoFillDelete(s, tx) }
```

**app/resources/Rs_Example.go** — gunakan `ToResource(c, data, fn)` wrapper
```go
func ExampleResource(c *fiber.Ctx, data any) any {
    return ToResource(c, data, ExampleSingleResource)
}
func ExampleSingleResource(c *fiber.Ctx, data any) ResponseExample {
    dataMap, _ := data.(map[string]any)
    return ResponseExample{
        Id:   helpers.Conv(dataMap["id"]).String(),
        Name: helpers.Conv(dataMap["name"]).String(),
    }
}
```

---

## File Upload (Cloudflare R2)

- Profile photo: upload langsung dari browser ke presigned URL
- Step images (future): sama, presigned URL dari API
- Format yang didukung: JPG, PNG, WebP. Max 5MB per file.

---

## Design Brief — "Cozy / Forest Green + Cream + Warm Tan"

> **Aktif**: tema Cozy diterapkan di codebase. Prototype asli (Italian Zine) tersimpan di `/tmp/cozy-extract/yava/project/YAVA Prototype.html` untuk referensi.

Handoff bundle: `/tmp/cozy-extract/yava/project/` — termasuk `YAVA Prototype - Cozy.html`, `styles-cozy.css`, `cozy-decorations.js`.

### Filosofi

Café reading-nook hangat. Forest green sebagai primer, cream sebagai latar, warm tan sebagai aksi. Serif hangat (Spectral) untuk headlines, Manrope untuk UI. Radius lebih besar, warm shadow, paper-grain halus di latar dan dark surface. Tidak brutal — homey.

### Palet Warna

| Token | Hex | Penggunaan |
|---|---|---|
| `--coral-red` | `#AD8257` | CTA utama, aksen, dot, highlight (warm tan) |
| `--coral-deep` | `#8C6741` | Hover state primary button |
| `--electric` | `#2D524A` | Grup/komunitas, card electric (forest green) |
| `--powder` | `#E1BF91` | Timer numerals on dark surfaces |
| `--lilac` | `#E1BF91` | Subtext on dark/electric surfaces (= powder) |
| `--deep-ink` | `#2D524A` | Teks utama, dark cards (forest green) |
| `--abyss` | `#1B342E` | Full-bleed brewing mode background |
| `--lavender-fog` | `#EEEEEC` | App background (cream, bukan putih!) |
| `--grid-paper` | `#CBC0AC` | Editorial brewing bg, tag grup (sage oat) |
| `--hairline` | `#DDD3BE` | Dividers, borders cards |
| `--muted` | `#8A8273` | Secondary text |
| `--white` (surface) | `#FBF8F1` | Card background (warm off-white) |

### Tipografi

- **Display/H1/H2**: Spectral (serif, italic untuk logo) — weight 600
- **UI**: Manrope — weight 400–700
- `.t-display`: 64px, Spectral 600, kerning -.025em
- `.t-h1` / `.t-h2`: 34px / 22px, Spectral 600
- `.t-h3`: 17px, Manrope 700
- `.t-label`: 11px, uppercase, letter-spacing .12em — kicker/kategori
- `.t-mono-num`: tabular-nums — semua angka timer dan stats
- **Logo**: Spectral italic, weight 600, letter-spacing -.03em

### Komponen Utama

**Topnav** — bukan sidebar. Logo kiri (Spectral italic), nav links tengah, search + bell + CTA + avatar kanan.

**Tags** — pill shape. Cozy: espresso=tan, v60=forest+lilac, cold=abyss, grup=grid-paper/sage.

**Cards** — `.card` (#FBF8F1, border-radius 18px), `.card--dark` (forest green, 22px, film grain overlay), `.card--electric` (abyss bg, lilac text). Dark surfaces punya `::after` grain pseudo-element — jangan hapus.

**Buttons** — pill. Primary=tan (#AD8257), light-primary=forest ink bg+lilac text, secondary=bordered forest.

**Params grid** — 5 kolom, background #FBF8F1.

**Step list** — `.step__num` pakai forest green bg + lilac text.

### Cozy Decorations

SVG dekorasi di `components/cozy-decorations.tsx`. Semua `position: absolute`, parent harus `position: relative; overflow: hidden`.

| Komponen | Lokasi |
|---|---|
| `<CozyFigureMug />` | `.card--dark.card--hero` di dashboard (sosok + mug) |
| `<CozyBranch />` | `.card--electric` (kartu grup) — ranting kopi + cherry |
| `<CozyMugSteam />` | `.card--dark.card--hero` di recipe detail (cangkir samping) |
| `<CozyPlants />` | Title row di recipe detail (monstera trio) |
| `<CozyFigureWalking />` | Coral complete block di editorial brewing mode |

### Brewing Mode — 3 Variants

| Variant | Background | Timer Size | Layout |
|---|---|---|---|
| `focus` | `--abyss` (#1B342E) | 260px | Centered, cream-tan timer on dark green |
| `ambient` | `--lavender-fog` | 200px | 2-col: dark card kiri, sidecar kanan |
| `editorial` | `--grid-paper` + grid lines | 360px | Asymmetric: text kiri, abyss block kanan |

Countdown: 3 → 2 → 1 dengan `.count-pop` 800ms. Timer numerals pakai `--powder` (#E1BF91).

### Aturan Tambahan

1. **Warna latar selalu** `--lavender-fog` (#EEEEEC) dengan paper-grain — BUKAN putih murni
2. **Logo** selalu Spectral italic: `YAVA<span class="dot">.</span>` — dot pakai `--coral-red` (#AD8257)
3. Angka timer selalu `t-mono-num` (tabular-nums)
4. Brewing mode **fullscreen** (`position: fixed; inset: 0`)
5. Dark cards (`.card--dark`, `.card--abyss`) punya `::after` grain — anak langsung harus `position: relative; z-index: 1` agar tidak tertutup grain (sudah di-handle oleh `> *` selector di globals.css)
6. Untuk komponen baru, cek dulu class yang ada di `globals.css` sebelum membuat style baru
