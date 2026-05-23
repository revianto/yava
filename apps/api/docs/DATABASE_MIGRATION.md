# Database Migration System

Panduan penggunaan migration system berbasis `golang-migrate` untuk mengelola skema database.

---

## File Structure

```
cmd/migrate/main.go         # CLI entry point
database/migrations/         # SQL migration files
```

---

## Commands

### Create Migration
```bash
go run cmd/migrate/main.go create <migration_name>
```

**Contoh:**
```bash
go run cmd/migrate/main.go create create_users_table
# Output: database/migrations/20260124140000_create_users_table.up.sql
```

### Run Migrations (Up)
```bash
go run cmd/migrate/main.go up
```

### Rollback Migrations (Down)
```bash
go run cmd/migrate/main.go down
```

### Force Version
```bash
go run cmd/migrate/main.go force <version_number>
```

---

## Migration Files

Setiap migration terdiri dari 2 file:

| File | Description |
|------|-------------|
| `*.up.sql` | SQL untuk apply perubahan (CREATE, ALTER, etc) |
| `*.down.sql` | SQL untuk revert perubahan (DROP, etc) |

**Contoh `up.sql`:**
```sql
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Contoh `down.sql`:**
```sql
DROP TABLE IF EXISTS users;
```

---

## Error Handling (Dirty State)

Jika migration gagal, database masuk ke **dirty state**. Untuk fix:

1. Perbaiki error di SQL file
2. Force ke versi terakhir yang sukses:
   ```bash
   go run cmd/migrate/main.go force <version_number>
   ```
3. Jalankan `up` lagi

---

## Technical Details

- Library: `golang-migrate/migrate/v4`
- Config: dari `.env` via `helper.GetEnv`
- Tracking: tabel `schema_migrations`

> [!IMPORTANT]
> Jalankan `go mod tidy` sebelum menggunakan commands di atas.
