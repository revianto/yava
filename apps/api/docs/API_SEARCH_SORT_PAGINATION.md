# Panduan Search, Sort & Pagination

Dokumentasi ini menjelaskan cara menggunakan fitur search, sort, dan pagination pada API endpoints yang didukung oleh Tornado Go Backend.

---

## Pagination

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `limit` | integer | 15 | Jumlah item per halaman |

### Contoh

```bash
# Halaman 1 dengan 10 items per halaman
GET /api/en/admin/item?page=1&limit=10
```

### Response Structure

Pagination metadata dikembalikan di dalam properti `meta`.

```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 45,
    "total_pages": 5
  }
}
```

---

## Sort

Fitur sort mendpaat input dari query parameters atau JSON body.

### Query Parameters (Single Sort)

| Parameter | Type | Description |
|-----------|------|-------------|
| `sort[field]` | string | Nama kolom untuk sorting |
| `sort[value]` | string | Arah sorting: `asc` (default) atau `desc` |

```bash
# Sort by name ascending
GET /api/en/admin/product?sort[field]=name&sort[value]=asc
```

### Query Parameters (Multiple/Indexed Sort)

Untuk sorting lebih dari satu kolom, gunakan index array (berurutan dari 0):

```bash
# Sort by name asc, then by id desc
GET /api/en/admin/product?sort[0][field]=name&sort[0][value]=asc&sort[1][field]=id&sort[1][value]=desc
```

### JSON Body Alternative

```json
{
  "sort": [
    { "field": "name", "value": "asc" },
    { "field": "base_price", "value": "desc" }
  ]
}
```

---

## Search

Search menggunakan format nested query parameters atau JSON body untuk mendukung kondisi kompleks (AND/OR).

### Single Condition

```bash
# Search where name contains "coffee"
GET /api/en/admin/product?search[field]=name&search[operator]=like&search[value]=coffee
```

### AND Conditions

Menggunakan array index `[0]`, `[1]`, dst:

```bash
# Search where name contains "coffee" AND unit_id = 5
GET /api/en/admin/product?search[and][0][field]=name&search[and][0][operator]=like&search[and][0][value]=coffee&search[and][1][field]=unit_id&search[and][1][operator]==&search[and][1][value]=5
```

### OR Conditions

```bash
# Search where name contains "coffee" OR name contains "tea"
GET /api/en/admin/product?search[or][0][field]=name&search[or][0][operator]=like&search[or][0][value]=coffee&search[or][1][field]=name&search[or][1][operator]=like&search[or][1][value]=tea
```

### Supported Operators

Operator yang diizinkan tergantung pada konfigurasi model. Operator umum yang didukung oleh sistem core:

| Operator | SQL Equivalent | Description |
|----------|----------------|-------------|
| `=` | `= ?` | Equal |
| `!=` | `!= ?` | Not equal |
| `>` | `> ?` | Greater than |
| `>=` | `>= ?` | Greater than or equal |
| `<` | `< ?` | Less than |
| `<=` | `<= ?` | Less than or equal |
| `like` | `LIKE %?%` | Contains (wildcard ditambahkan otomatis) |
| `is` | `IS ?` | Digunakan untuk `IS NULL` atau `IS NOT NULL` |

### JSON Body Alternative

Gunakan JSON body untuk kondisi yang sangat kompleks:

```json
{
  "search": {
    "and": [
      { "field": "name", "operator": "like", "value": "coffee" },
      {
        "or": [
          { "field": "unit_id", "operator": "=", "value": "1" },
          { "field": "unit_id", "operator": "=", "value": "2" }
        ]
      }
    ]
  }
}
```

---

## Searchable & Sortable Fields

Daftar field yang tersedia tergantung pada masing-masing model.

### Product
- **Searchable**: `id`, `name`, `unit_id`
- **Sortable**: `id`, `name`, `base_price`

### Customer
- **Searchable**: `id`, `company_name`, `owner_name`, `phone`, `email`
- **Sortable**: `id`, `company_name`, `total_orders`, `total_debt`

### Role
- **Searchable**: `id`, `code`, `name`
- **Sortable**: `id`, `name`

---

## Tips

1. **URL Encoding**: Selalu gunakan URL encoding untuk nilai query parameter dengan spasi atau karakter khusus.
2. **Sequential Index**: Saat menggunakan `sort[n]` atau `search[and][n]`, index `n` harus berurutan mulai dari 0.
3. **Field Names**: Pastikan nama field sesuai dengan yang didefinisikan di interface model (alias pada query select).
4. **Case Sensitivity**: Operator string seperti `asc`, `desc`, dan `like` tidak case-sensitive di level parsing (diubah ke lowercase/uppercase otomatis).

## Contoh Penggunaan di Postman

### 1. Pagination Saja
- **URL**: `http://localhost:3000/api/en/admin/product?page=2&limit=10`

### 2. Sort Saja (Multiple)
- **URL**: `http://localhost:3000/api/en/admin/product?sort[0][field]=name&sort[0][value]=asc&sort[1][field]=base_price&sort[1][value]=desc`

### 3. Search Saja (AND + OR)
- **URL**: `http://localhost:3000/api/en/admin/product?search[and][0][field]=name&search[and][0][operator]=like&search[and][0][value]=coffee&search[and][1][or][0][field]=unit_id&search[and][1][or][0][operator]==&search[and][1][or][0][value]=1&search[and][1][or][1][field]=unit_id&search[and][1][or][1][operator]==&search[and][1][or][1][value]=2`

### 4. Kombinasi Lengkap
- **URL**: `http://localhost:3000/api/en/admin/product?page=2&limit=10&sort[0][field]=name&sort[0][value]=asc&search[and][0][field]=name&search[and][0][operator]=like&search[and][0][value]=coffee`