# Database Migrations Guide

Proyek ini menggunakan [golang-migrate](https://github.com/golang-migrate/migrate) untuk manajemen skema database.

## Prasyarat

Anda perlu menginstal CLI `migrate`.

### Menggunakan Homebrew (macOS)
```bash
brew install golang-migrate
```

### Menggunakan Go Install
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Menggunakan Docker
Jika Anda tidak ingin menginstal binary lokal, Anda bisa menyesuaikan Makefile untuk menggunakan `docker run`.

## Konfigurasi

Pastikan file `.env` Anda memiliki variabel `POSTGRES_DSN` yang valid.
Contoh:
```env
POSTGRES_DSN=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

## Perintah Migration (Makefile)

Gunakan `make` untuk menjalankan perintah migrasi dengan mudah.

| Perintah | Deskripsi | Contoh |
|----------|-----------|--------|
| `make migrate-create` | Membuat file migrasi baru (up & down) | `make migrate-create name=add_users_table` |
| `make migrate-up` | Menjalankan semua migrasi yang belum terapkan (UP) | `make migrate-up` |
| `make migrate-down` | Membatalkan migrasi terakhir (DOWN 1 langkah) | `make migrate-down` |
| `make migrate-force` | Memaksa database ke versi tertentu (berguna jika dirty) | `make migrate-force version=20240520120000` |
| `make migrate-version` | Cek versi migrasi saat ini | `make migrate-version` |

## Struktur Migrasi

File migrasi disimpan di folder `migrations/`. Setiap migrasi terdiri dari dua file:
- `YYYYMMDDHHMMSS_name.up.sql`: Perubahan yang akan diterapkan.
- `YYYYMMDDHHMMSS_name.down.sql`: Cara membatalkan perubahan tersebut (rollback).

## Error Handling & Troubleshooting

### Dirty Database
Jika migrasi gagal di tengah jalan, database akan ditandai sebagai "dirty". Anda perlu memperbaikinya secara manual (misalnya menghapus tabel yang setengah jadi) lalu memaksa versi ke versi terakhir yang sukses.

```bash
# Contoh: Force ke versi sebelumnya agar bisa di-run ulang
make migrate-force version=20240520110000
```

### Driver Database
Setup ini dikonfigurasi untuk PostgreSQL (`postgres://` atau `postgresql://`). Jika menggunakan MySQL, ubah driver di URL koneksi menjadi `mysql://`.

## Development vs Production

### Development
Di local, gunakan `make migrate-up` untuk selalu sync dengan skema terbaru.

### Production
Di production, disarankan menjalankan migrasi sebagai bagian dari pipeline deployment (CI/CD) atau menggunakan init container (jika di Kubernetes), bukan menjalankan `migrate-up` secara manual dari laptop developer.

Pastikan environment variable `POSTGRES_DSN` di production mengarah ke database production dengan user yang memiliki hak akses DDL.
