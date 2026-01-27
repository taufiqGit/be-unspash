# REST API MVC Sederhana (Go)

Proyek ini membuat REST API sederhana dengan arsitektur MVC untuk resource `Todo` menggunakan standard library `net/http` (tanpa framework eksternal).

## Struktur Proyek

```
.
├── app.go                # entrypoint server
├── controllers/          # handler HTTP
├── models/               # definisi model/data
├── routes/               # registrasi routing
├── services/             # logika bisnis + storage in-memory
└── go.mod
```

## Menjalankan

1. Pastikan Go terinstall (disarankan Go 1.21+).
2. Jalankan server:

```bash
go run .
```

Server akan berjalan di `http://localhost:8080`.

## Endpoint

- `GET /api/todos` — ambil daftar todo
- `POST /api/todos` — buat todo baru
  - Body JSON: `{ "title": "Belajar Go", "done": false, "image_url": "https://example.com/img.png" }`
- `GET /api/todos/{id}` — ambil todo by ID
- `PUT /api/todos/{id}` — update todo by ID
  - Body JSON: `{ "title": "Belajar Go Lanjut", "done": true, "image_url": "https://example.com/new.png" }`
- `DELETE /api/todos/{id}` — hapus todo by ID

### Categories

- `GET /api/categories` — ambil semua kategori
- `POST /api/categories` — buat kategori baru
  - Body JSON: `{ "name": "Teknologi", "description": "Kategori seputar teknologi" }`
- `GET /api/categories/{id}` — ambil detail kategori
- `PUT /api/categories/{id}` — update kategori
  - Body JSON: `{ "name": "Teknologi Baru", "description": "Update deskripsi" }`
- `DELETE /api/categories/{id}` — hapus kategori

## Format Respons

Semua respons menggunakan envelope konsisten:

```json
{
  "success": true,
  "message": "list todos",
  "data": [ { "id": 1, "title": "Belajar Go", "done": false, "image_url": "https://example.com/img.png", "created_at": "...", "updated_at": "..." } ],
  "meta": { "count": 1 }
}
```

Contoh error:

```json
{
  "success": false,
  "error": { "code": "VALIDATION_ERROR", "message": "title dan image url tidak boleh kosong" }
}
```

## Contoh uji dengan curl

```bash
# Buat todo baru
curl -s -X POST http://localhost:8080/api/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"Belajar Go","done":false, "image_url":"https://example.com/img.png"}' | jq

# List todos
curl -s http://localhost:8080/api/todos | jq

# Ambil todo id=1
curl -s http://localhost:8080/api/todos/1 | jq

# Update todo id=1
curl -s -X PUT http://localhost:8080/api/todos/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Belajar Go Lanjut","done":true, "image_url":"https://example.com/new.png"}' | jq

# Hapus todo id=1
curl -s -X DELETE http://localhost:8080/api/todos/1 | jq
```

Catatan: data disimpan in-memory, sehingga akan hilang saat server dimatikan.

## Docker

Build image:

```bash
docker build -t gowes-api .
```

Jalankan container:

```bash
docker run -d --name nama_container -p 8080:8080 gowes-api
```

Server akan tersedia di `http://localhost:8080`.
