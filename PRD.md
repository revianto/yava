# YAVA — Product Requirements Document

> **Version:** 1.0
> **Platform:** Web Application
> **Date:** Mei 2026
> **Status:** Pre-Development

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Tech Stack](#2-tech-stack)
3. [User Roles & Personas](#3-user-roles--personas)
4. [Information Architecture](#4-information-architecture)
5. [Data Models](#5-data-models)
6. [Features](#6-features)
7. [Timer System](#7-timer-system)
8. [Recipe Visibility & Sharing](#8-recipe-visibility--sharing)
9. [Group System](#9-group-system)
10. [User Flows](#10-user-flows)
11. [API Contract](#11-api-contract)
12. [Non-Functional Requirements](#12-non-functional-requirements)
13. [Development Roadmap](#13-development-roadmap)
14. [Future Scope](#14-future-scope)
15. [Glossary](#15-glossary)

---

## 1. Project Overview

### 1.1 What is YAVA?

YAVA (**Y**our **A**mazing **V**arious **A**romas) adalah web application untuk menyimpan, mengelola, dan berbagi resep kopi.

**Core Value:**
- Resep terstruktur dengan hierarki `Jenis → Sub-jenis → Resep`
- Timer otomatis multi-sesi yang mengalir **tanpa jeda** antar tahap brewing
- Kontrol privasi resep: Pribadi / Publik / Grup
- Komunitas kopi berbasis grup dengan sistem peran

### 1.2 Scope v1.0

| Kategori | Included (v1.0) | Excluded (Future) |
|---|---|---|
| Platform | Web (Next.js) | iOS, Android |
| Auth | Google OAuth / SSO | Email/password |
| Resep | CRUD + 3-level visibility | Recipe versioning |
| Timer | Multi-session auto-advance + 3s countdown | Custom audio cues |
| Grup | Admin/member roles, recipe approval, diskusi | Sub-groups |
| AI | — | AI recipe recommendation |

---

## 2. Tech Stack

| Layer | Technology | Notes |
|---|---|---|
| Frontend | Next.js (App Router) | SSR + CSR hybrid |
| Backend | Golang (REST API) | Stateless, Docker-ready |
| Database | PostgreSQL | Relational |
| Auth | Google OAuth 2.0 | JWT session token |
| File Storage | S3-compatible (Cloudflare R2) | Profile photos, step images |
| Realtime | WebSocket / SSE | In-app notifications |
| Frontend Deploy | Vercel | |
| Backend Deploy | VPS + Docker | |

---

## 3. User Roles & Personas

### 3.1 Application-Level Roles

| Role | Description |
|---|---|
| `user` | Akun pengguna standar |
| `system_admin` | Pengelola resep default (internal YAVA) |

### 3.2 Group-Level Roles

| Role | Description |
|---|---|
| `group_admin` | Dapat approve resep, kelola member, hapus grup |
| `group_member` | Dapat share resep, lihat koleksi, ikut diskusi |

### 3.3 Persona Summary

| Persona | Background | Primary Need |
|---|---|---|
| Home Brewer | Baru mulai brewing | Timer otomatis, resep default sebagai referensi |
| Barista Profesional | Brewing harian, resep kompleks | Manajemen resep detail, share ke tim |
| Komunitas Admin | Mengelola grup kopi | Control approve resep, manajemen member |

---

## 4. Information Architecture

### 4.1 Recipe Hierarchy

```
Jenis (Type)
└── Sub-jenis (Subtype)
    └── Resep (Recipe)
        ├── Parameters (dose, yield, temp, grind, ratio)
        ├── Sessions[] (timer-based steps)
        └── Notes[] (non-timer steps, ordered)
```

### 4.2 Example Hierarchy

```
V60
└── Regular Drip
    └── "V60 Light Roast 15g/250ml"
        ├── Params: 15g kopi, 250ml air, 92°C, Medium-Fine
        ├── Session 1: Blooming — 45s
        ├── Session 2: First Pour — 100s
        └── Session 3: Second Pour — 100s

Espresso
├── Machine
│   └── "Espresso Classic 18g/36g"
└── Manual (Flair)
    └── "Flair Espresso 18g/36g 9Bar"
        ├── Session 1: Pre-infusion — 30s
        ├── Session 2: Extraction — 30s
        └── Note: "Campur dengan susu steamed 150ml, aduk perlahan"
```

---

## 5. Data Models

### 5.1 User

```go
type User struct {
    ID          uuid.UUID `db:"id"`
    GoogleID    string    `db:"google_id"`
    Email       string    `db:"email"`
    Name        string    `db:"name"`
    AvatarURL   string    `db:"avatar_url"`
    CreatedAt   time.Time `db:"created_at"`
}
```

### 5.2 RecipeType (Jenis)

```go
type RecipeType struct {
    ID        uuid.UUID `db:"id"`
    Name      string    `db:"name"`      // "V60", "Espresso"
    IsDefault bool      `db:"is_default"`
}
```

### 5.3 RecipeSubtype (Sub-jenis)

```go
type RecipeSubtype struct {
    ID         uuid.UUID `db:"id"`
    TypeID     uuid.UUID `db:"type_id"`
    Name       string    `db:"name"`     // "Manual (Flair)", "Machine"
    IsDefault  bool      `db:"is_default"`
}
```

### 5.4 Recipe

```go
type Recipe struct {
    ID          uuid.UUID      `db:"id"`
    OwnerID     *uuid.UUID     `db:"owner_id"`   // NULL = system/default recipe
    TypeID      uuid.UUID      `db:"type_id"`
    SubtypeID   uuid.UUID      `db:"subtype_id"`
    Name        string         `db:"name"`
    Description string         `db:"description"`

    // Brewing parameters
    CoffeeDoseGram  float32 `db:"coffee_dose_gram"`
    YieldGram       float32 `db:"yield_gram"`
    WaterTempCelsius float32 `db:"water_temp_celsius"`
    GrindSize       string  `db:"grind_size"`   // "Fine" | "Medium-Fine" | "Medium" | "Coarse"
    Ratio           string  `db:"ratio"`        // e.g. "1:15"

    // Visibility & status
    Visibility  string    `db:"visibility"`   // "private" | "public" | "group"
    Status      string    `db:"status"`       // "active" | "archived"
    IsDefault   bool      `db:"is_default"`   // system recipe, cannot be deleted

    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}
```

### 5.5 RecipeSession (Sesi Timer)

```go
type RecipeSession struct {
    ID           uuid.UUID `db:"id"`
    RecipeID     uuid.UUID `db:"recipe_id"`
    Order        int       `db:"order"`          // 1-based, determines execution sequence
    Name         string    `db:"name"`           // "Blooming", "First Pour"
    DurationSec  int       `db:"duration_sec"`   // seconds
    GuideNote    string    `db:"guide_note"`     // optional instruction shown during session
}
```

### 5.6 RecipeNote (Notes Non-Timer)

```go
type RecipeNote struct {
    ID         uuid.UUID `db:"id"`
    RecipeID   uuid.UUID `db:"recipe_id"`
    Order      int       `db:"order"`       // determines display position among sessions + notes
    Content    string    `db:"content"`     // free text
}
```

> **Note on ordering:** `RecipeSession.Order` and `RecipeNote.Order` share the same sequence space. Example: Order 1 = Note, Order 2 = Session (Blooming), Order 3 = Session (First Pour), Order 4 = Note.

### 5.7 Group

```go
type Group struct {
    ID          uuid.UUID `db:"id"`
    Name        string    `db:"name"`
    Description string    `db:"description"`
    FounderID   uuid.UUID `db:"founder_id"`
    InviteCode  string    `db:"invite_code"`   // unique link for joining
    CreatedAt   time.Time `db:"created_at"`
}
```

### 5.8 GroupMember

```go
type GroupMember struct {
    GroupID   uuid.UUID `db:"group_id"`
    UserID    uuid.UUID `db:"user_id"`
    Role      string    `db:"role"`       // "admin" | "member"
    JoinedAt  time.Time `db:"joined_at"`
}
```

### 5.9 GroupRecipe

```go
type GroupRecipe struct {
    GroupID    uuid.UUID  `db:"group_id"`
    RecipeID   uuid.UUID  `db:"recipe_id"`
    Status     string     `db:"status"`        // "pending" | "active" | "archived" | "rejected"
    SubmittedBy uuid.UUID `db:"submitted_by"`
    ReviewedBy  *uuid.UUID `db:"reviewed_by"`  // admin who approved/rejected
    SubmittedAt time.Time  `db:"submitted_at"`
    ReviewedAt  *time.Time `db:"reviewed_at"`
}
```

### 5.10 Discussion

```go
type Discussion struct {
    ID        uuid.UUID  `db:"id"`
    GroupID   uuid.UUID  `db:"group_id"`
    RecipeID  uuid.UUID  `db:"recipe_id"`
    UserID    uuid.UUID  `db:"user_id"`
    ParentID  *uuid.UUID `db:"parent_id"`  // NULL = top-level comment, set = reply
    Content   string     `db:"content"`
    IsPinned  bool       `db:"is_pinned"`  // admin-only
    CreatedAt time.Time  `db:"created_at"`
}
```

### 5.11 Notification

```go
type Notification struct {
    ID        uuid.UUID `db:"id"`
    UserID    uuid.UUID `db:"user_id"`
    Type      string    `db:"type"`    // "recipe_approved" | "recipe_rejected" | "group_invite" | "discussion_reply"
    Payload   jsonb     `db:"payload"` // context-specific JSON
    IsRead    bool      `db:"is_read"`
    CreatedAt time.Time `db:"created_at"`
}
```

---

## 6. Features

### 6.1 Authentication

- Login via **Google OAuth 2.0** only (v1.0)
- Successful OAuth → server issues **JWT** (stored in HttpOnly cookie)
- One Google account = one YAVA profile
- No manual email/password registration

### 6.2 Recipe Management

#### Default Recipes (System)
- Created and managed by `system_admin`
- `owner_id = NULL`, `is_default = true`
- Visible to **all users** automatically
- **Cannot be deleted or edited** by regular users
- Users can **duplicate** a default recipe to their own collection → creates new recipe with `owner_id = user.id`, `is_default = false`

#### Private Recipes
- `visibility = "private"`, only visible to owner
- Full CRUD by owner
- Can be promoted to `public` or submitted to a group

#### Public Recipes
- `visibility = "public"`, visible to all users
- **Only owner can edit**
- Owner can archive → `status = "archived"` (not deleted, hidden from public search)
- Other users can duplicate to their private collection

#### Group Recipes
- Recipe submitted to a group, requires **admin approval** before activation
- `status` flow: `pending` → `active` (approved) or `rejected`
- Only owner can edit recipe content
- Group admin can remove recipe from group → `status = "archived"` in `GroupRecipe`
- Owner can withdraw recipe from group → `status = "archived"` in `GroupRecipe`
- Original recipe in owner's collection is **never deleted**

### 6.3 Recipe CRUD Rules

| Action | Who |
|---|---|
| Create | Any authenticated user |
| Read | Owner (private), All users (public/default), Group members (group) |
| Update | Owner only |
| Delete | Owner only (soft delete / archive, not hard delete) |
| Duplicate | Any user with read access |

---

## 7. Timer System

### 7.1 Behavior Specification

```
User clicks "Mulai Brewing"
        │
        ▼
┌─────────────────────┐
│  Countdown 3 detik  │  ← Preparation only, NOT a session
└─────────────────────┘
        │ auto-advance
        ▼
┌─────────────────────┐
│    Session 1        │  ← Timer counts down
│    e.g. Blooming    │
│    45s              │
└─────────────────────┘
        │ auto-advance (NO gap)
        ▼
┌─────────────────────┐
│    Session 2        │
│    e.g. First Pour  │
│    100s             │
└─────────────────────┘
        │ auto-advance (NO gap)
        ▼
┌─────────────────────┐
│    Session 3        │
│    e.g. Second Pour │
│    100s             │
└─────────────────────┘
        │ last session ends
        ▼
   Brewing Complete Screen
```

**Key rules:**
- Transition between sessions: **zero delay**, immediate
- The 3-second preparation countdown is **outside** the session sequence
- Sessions execute in ascending `order` value
- `RecipeNote` items are displayed at their `order` position but **do not block** the timer

### 7.2 Timer Controls

| Control | Behavior |
|---|---|
| Pause | Freezes current session countdown |
| Resume | Continues from paused time |
| Reset | Returns to 3s preparation countdown |
| Skip Session | Immediately advances to next session |

### 7.3 Timer Display

```
┌─────────────────────────────────────────┐
│  Blooming                               │
│                                         │
│            00:32                        │  ← remaining time
│                                         │
│  Session 1 of 3                         │
│  ████████░░░░░░░░░░  (progress bar)     │
│                                         │
│  [  Pause  ]    [  Skip  ]             │
└─────────────────────────────────────────┘
```

### 7.4 Recipes Without Timer

- A recipe **may have zero sessions** — timer feature is optional
- Notes-only recipes display steps as a checklist (no countdown)
- Example: Kopi Susu recipe may have timer sessions for espresso extraction but plain notes for milk mixing steps

---

## 8. Recipe Visibility & Sharing

### 8.1 Visibility States

```
private ──► public    (owner promotes)
private ──► group     (owner submits, needs admin approval)
public  ──► archived  (owner archives)
group   ──► archived  (owner withdraws, or group admin removes)
archived ◄──► active  (owner can restore)
```

### 8.2 Visibility Matrix

| Who can see | private | public | group | archived |
|---|---|---|---|---|
| Owner | ✓ | ✓ | ✓ | ✓ (labeled) |
| All users | ✗ | ✓ | ✗ | ✗ |
| Group members | ✗ | ✗ | ✓ (if active) | ✗ |

### 8.3 Archive Behavior

- Archived recipes are **NOT deleted**
- Removed from public search and group collections
- Still accessible by owner in their dashboard with "Archived" label
- Owner can restore to `active` at any time

---

## 9. Group System

### 9.1 Group Lifecycle

```
User creates group → becomes founder (auto: group_admin)
        │
        ▼
Invite members via invite_code link or email search
        │
        ▼
Members join → role: "member"
        │
        ▼
Admin can promote member → role: "admin"
```

### 9.2 Permission Matrix

| Action | group_admin | group_member |
|---|---|---|
| Approve / reject submitted recipes | ✓ | ✗ |
| Remove recipe from group | ✓ | ✗ |
| Submit recipe to group | ✓ | ✓ |
| View & use group recipes | ✓ | ✓ |
| Post comments / discussions | ✓ | ✓ |
| Pin a comment | ✓ | ✗ |
| Invite new members | ✓ | ✗ |
| Promote member to admin | ✓ | ✗ |
| Delete / disband group | ✓ (founder only) | ✗ |

### 9.3 Recipe Submission Flow

```
Member submits recipe to group
        │
        ▼
GroupRecipe.status = "pending"
        │
        ▼
Admin receives notification
        │
   ┌────┴────┐
   ▼         ▼
Approve    Reject
   │         │
   ▼         ▼
status=    status=      Owner receives notification
"active"  "rejected"    (with optional rejection reason)
```

### 9.4 Discussions

- Each recipe in a group has its own discussion thread
- Comments support nested replies (`parent_id`)
- Admin can pin (`is_pinned = true`) important comments
- User receives notification when their comment gets a reply

---

## 10. User Flows

### 10.1 Onboarding

```
Visit YAVA → Click "Login with Google"
        │
        ▼
Google OAuth consent → callback to /auth/google
        │
        ▼
Server: create or fetch User record → issue JWT
        │
        ▼
Redirect to /dashboard
```

### 10.2 Create New Recipe

```
/dashboard → "Buat Resep"
        │
        ▼
Step 1: Fill basic info
  - Name, description
  - Select Type (Jenis)
  - Select Subtype (Sub-jenis)
  - Brewing parameters: dose, yield, temp, grind, ratio
        │
        ▼
Step 2: Add Sessions (optional, repeatable)
  - Session name
  - Duration (seconds)
  - Guide note
  - Drag to reorder
        │
        ▼
Step 3: Add Notes (optional, repeatable)
  - Free text content
  - Set position (order) among sessions
        │
        ▼
Step 4: Set Visibility
  - Private / Public / Group
  - If Group: select target group(s)
        │
        ▼
Save → POST /recipes
```

### 10.3 Share Recipe to Group

```
Recipe detail page → "Bagikan ke Grup"
        │
        ▼
Select group (user must be a member)
        │
        ▼
POST /groups/:id/recipes
        │
        ▼
GroupRecipe.status = "pending"
Admin notified
        │
   ┌────┴────┐
   ▼         ▼
Approved   Rejected
   │         │
Owner      Owner
notified   notified + reason
```

### 10.4 Archive Recipe

```
Recipe detail → kebab menu → "Archive"
        │
        ▼
PATCH /recipes/:id/archive
        │
        ▼
status = "archived"
  - Removed from public search
  - Removed from group collections (GroupRecipe.status = "archived")
  - Still visible to owner in dashboard (labeled "Archived")
  - Restorable: PATCH /recipes/:id/restore
```

### 10.5 Brewing Session

```
Recipe detail page → "Mulai Brewing"
        │
        ▼
Full-screen brewing mode
        │
        ▼
[3-second preparation countdown]
  - Display: "Siapkan peralatan Anda..."
  - NOT counted as a session
        │
        ▼
Auto-start Session 1
  - Display: session name, countdown, progress bar
  - Available controls: Pause, Skip
        │ (session ends, 0 delay)
        ▼
Auto-start Session 2
  ...
        │ (last session ends)
        ▼
Brewing Complete Screen
  - Summary: total time, sessions completed
  - Actions: "Ulangi", "Kembali ke Resep"
```

---

## 11. API Contract

### 11.1 Base URL

```
Production : https://api.yava.app/v1
Development: http://localhost:8080/v1
```

### 11.2 Authentication

All endpoints (except `/auth/*`) require:

```
Authorization: Bearer <jwt_token>
```

JWT payload:
```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "exp": 1234567890
}
```

### 11.3 Endpoints

#### Auth

| Method | Path | Description |
|---|---|---|
| `GET` | `/auth/google` | Redirect to Google OAuth consent |
| `GET` | `/auth/google/callback` | OAuth callback, issue JWT |
| `POST` | `/auth/logout` | Invalidate session |

#### Recipes

| Method | Path | Description | Auth |
|---|---|---|---|
| `GET` | `/recipes` | List recipes. Query: `?visibility=private\|public`, `?type_id=`, `?page=`, `?limit=` | Required |
| `POST` | `/recipes` | Create new recipe | Required |
| `GET` | `/recipes/:id` | Get recipe detail (includes sessions + notes ordered) | Required |
| `PUT` | `/recipes/:id` | Update recipe (owner only) | Required |
| `PATCH` | `/recipes/:id/archive` | Archive recipe | Required |
| `PATCH` | `/recipes/:id/restore` | Restore archived recipe | Required |
| `POST` | `/recipes/:id/duplicate` | Duplicate recipe to caller's private collection | Required |

#### Recipe Types & Subtypes

| Method | Path | Description |
|---|---|---|
| `GET` | `/types` | List all recipe types |
| `GET` | `/types/:id/subtypes` | List subtypes for a type |

#### Groups

| Method | Path | Description | Auth |
|---|---|---|---|
| `POST` | `/groups` | Create group | Required |
| `GET` | `/groups/:id` | Group detail | Required (member) |
| `PUT` | `/groups/:id` | Update group info | Required (admin) |
| `DELETE` | `/groups/:id` | Disband group | Required (founder) |
| `GET` | `/groups/:id/members` | List members | Required (member) |
| `POST` | `/groups/:id/members` | Join via invite code | Required |
| `DELETE` | `/groups/:id/members/:uid` | Remove member | Required (admin) |
| `PATCH` | `/groups/:id/members/:uid/role` | Change member role | Required (admin) |

#### Group Recipes

| Method | Path | Description | Auth |
|---|---|---|---|
| `GET` | `/groups/:id/recipes` | List active group recipes | Required (member) |
| `GET` | `/groups/:id/recipes/pending` | List pending approval | Required (admin) |
| `POST` | `/groups/:id/recipes` | Submit recipe to group | Required (member) |
| `PATCH` | `/groups/:id/recipes/:rid/approve` | Approve recipe | Required (admin) |
| `PATCH` | `/groups/:id/recipes/:rid/reject` | Reject recipe (body: reason) | Required (admin) |
| `DELETE` | `/groups/:id/recipes/:rid` | Remove recipe from group | Required (admin) |

#### Discussions

| Method | Path | Description | Auth |
|---|---|---|---|
| `GET` | `/groups/:id/recipes/:rid/discussions` | Get discussion thread | Required (member) |
| `POST` | `/groups/:id/recipes/:rid/discussions` | Post comment | Required (member) |
| `PATCH` | `/groups/:id/recipes/:rid/discussions/:did/pin` | Pin comment | Required (admin) |
| `DELETE` | `/groups/:id/recipes/:rid/discussions/:did` | Delete comment (own or admin) | Required |

#### Notifications

| Method | Path | Description | Auth |
|---|---|---|---|
| `GET` | `/notifications` | List user notifications | Required |
| `PATCH` | `/notifications/:id/read` | Mark as read | Required |
| `PATCH` | `/notifications/read-all` | Mark all as read | Required |

### 11.4 Standard Response Format

**Success:**
```json
{
  "success": true,
  "data": { },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": "RECIPE_NOT_FOUND",
    "message": "Recipe with id xxx does not exist"
  }
}
```

### 11.5 Error Codes

| Code | HTTP Status | Meaning |
|---|---|---|
| `UNAUTHORIZED` | 401 | Missing or invalid JWT |
| `FORBIDDEN` | 403 | Action not allowed for this user/role |
| `NOT_FOUND` | 404 | Resource does not exist |
| `VALIDATION_ERROR` | 422 | Invalid request body |
| `CONFLICT` | 409 | e.g. recipe already submitted to group |
| `RATE_LIMITED` | 429 | Too many requests |

---

## 12. Non-Functional Requirements

### 12.1 Performance

| Metric | Target |
|---|---|
| Page load (recipe list, 4G) | < 2 seconds |
| Timer accuracy | ±100ms (use `performance.now()`) |
| API response time (CRUD) | < 500ms |
| API response time (list with pagination) | < 1 second |

### 12.2 Security

- All API endpoints require valid JWT (except `/auth/*`)
- **Server-side** ownership validation on every mutating request — do not rely on frontend
- Rate limiting:
  - `/auth/*`: 10 req/min per IP
  - `POST /recipes`: 30 req/min per user
  - `POST /groups/:id/recipes`: 20 req/min per user
- Input sanitization to prevent XSS in `content`, `guide_note`, `description` fields
- HttpOnly + Secure cookie for JWT storage

### 12.3 Scalability

- Stateless backend — horizontally scalable
- Mandatory pagination on all list endpoints (default: `limit=20`, max: `100`)
- Required database indexes:

```sql
CREATE INDEX idx_recipes_owner_id      ON recipes(owner_id);
CREATE INDEX idx_recipes_visibility    ON recipes(visibility, status);
CREATE INDEX idx_recipes_type_id       ON recipes(type_id);
CREATE INDEX idx_group_members_user    ON group_members(user_id);
CREATE INDEX idx_group_recipes_status  ON group_recipes(group_id, status);
CREATE INDEX idx_discussions_recipe    ON discussions(recipe_id, group_id);
CREATE INDEX idx_notifications_user    ON notifications(user_id, is_read);
```

---

## 13. Development Roadmap

| Phase | Duration | Deliverables |
|---|---|---|
| **Phase 1** | Week 1–2 | Project setup (Next.js + Golang), Google OAuth, DB schema, Recipe CRUD dasar |
| **Phase 2** | Week 3–4 | Timer multi-session auto-advance, notes non-timer, resep default, brewing mode UI |
| **Phase 3** | Week 5–6 | Recipe visibility (public/private), explore page, recipe duplication |
| **Phase 4** | Week 7–9 | Group system: CRUD, invite, recipe submission + approval, admin/member roles |
| **Phase 5** | Week 10–11 | Group discussions, in-app notifications, archive/restore |
| **Phase 6** | Week 12 | QA, performance testing, bug fixes, staging → production deploy |

---

## 14. Future Scope

- **AI Integration:** Recipe recommendation based on brewing profile
- **Mobile App:** Flutter (Android + iOS)
- **Recipe Versioning:** Change history + rollback
- **Rating & Review:** Public recipe ratings
- **Import / Export:** JSON format
- **Audio Timer:** Text-to-speech guidance per session

---

## 15. Glossary

| Term | Definition |
|---|---|
| **Sesi Timer (Session)** | Satu tahap brewing dengan durasi tertentu. Dieksekusi berurutan tanpa jeda. |
| **Notes Non-Timer** | Catatan panduan tanpa durasi, untuk langkah yang tidak memerlukan timer. |
| **Resep Default** | Resep bawaan sistem, tersedia untuk semua user, tidak bisa dihapus. |
| **Archived** | Status resep yang ditarik dari tampilan publik/grup tapi tetap ada di akun owner. |
| **Pending Approval** | Status resep yang disubmit ke grup, menunggu persetujuan admin. |
| **Countdown Persiapan** | Jeda 3 detik sebelum sesi pertama dimulai. Bukan bagian dari sesi brewing. |
| **Sub-jenis (Subtype)** | Turunan dari jenis resep berdasarkan peralatan/teknik. Contoh: Espresso → Flair Manual. |
| **Founder** | User yang membuat grup, secara default menjadi group_admin. Satu-satunya yang bisa membubarkan grup. |
| **Invite Code** | Kode unik per grup untuk undangan member baru. |
| **Auto-advance** | Perpindahan otomatis dari satu sesi ke sesi berikutnya tanpa interaksi pengguna. |
