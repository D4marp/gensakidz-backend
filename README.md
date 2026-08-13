# GenSA Kidz — Backend & Dashboard

Backend Go sederhana + dashboard admin untuk mengelola semua konten yang saat
ini tampil statis di gensakidz.com (layanan, karir, artikel, fasilitas,
galeri, tim, cabang, testimoni, dan teks umum seperti visi/misi). Sekali
dijalankan, seluruh konten yang sudah ada di situs otomatis dimasukkan ke
database — dashboard tidak dimulai dari kosong.

Satu binary, satu file database (SQLite), tanpa dependensi eksternal
(tidak perlu instal database server terpisah).

## Menjalankan

```bash
cd backend
go run .
```

Server jalan di `http://localhost:8080`, otomatis redirect ke `/admin`.

Login default (**wajib diganti** sebelum dipakai serius — lihat bagian
Konfigurasi di bawah):
- Email: `admin@gensakidz.com`
- Kata sandi: `gensakidz2026`

## Struktur

- `main.go` — routing & start server
- `db.go` — koneksi & skema database (SQLite, dibuat otomatis di `gensakidz.db`)
- `models.go` — struct data (Service, Job, Article, Facility, dst.)
- `repo.go` — query database per jenis konten
- `seed.go` — data awal (migrasi dari konten statis situs saat ini)
- `handlers.go` — halaman dashboard admin (list, tambah, ubah, hapus)
- `api.go` — API publik JSON read-only (`/api/*`) untuk dikonsumsi frontend
- `auth.go` — login sederhana (cookie session + bcrypt)
- `upload.go` — upload foto, disimpan di folder `uploads/`
- `templates/`, `static/` — tampilan dashboard (HTML + CSS polos, tanpa build step)

## Konfigurasi

Semua opsional, lewat environment variable:

| Variabel | Default | Keterangan |
|---|---|---|
| `PORT` | `8080` | Port server |
| `DB_PATH` | `gensakidz.db` | Lokasi file database |
| `ADMIN_EMAIL` | `admin@gensakidz.com` | Email login admin (dipakai hanya saat pertama kali dijalankan) |
| `ADMIN_PASSWORD` | `gensakidz2026` | Kata sandi admin (dipakai hanya saat pertama kali dijalankan) |

Contoh menjalankan dengan kredensial sendiri sejak awal:

```bash
ADMIN_EMAIL=owner@gensakidz.com ADMIN_PASSWORD="kata-sandi-kuat" go run .
```

> Kredensial admin hanya dibuat **sekali**, saat tabel `admin_users` masih
> kosong (database baru). Untuk mengganti kata sandi setelahnya, hapus baris
> di tabel `admin_users` lalu jalankan ulang dengan env var baru, atau
> tambahkan fitur ganti-password kalau nanti dibutuhkan.

## Apa yang belum otomatis

- **Foto**: data seed hanya berisi teks (judul, deskripsi, dll). Foto yang
  sudah ada di situs (Next.js) tidak ikut disalin otomatis ke sini karena
  keduanya adalah proyek terpisah. Upload ulang foto yang relevan lewat
  masing-masing form di dashboard (Layanan, Fasilitas, Galeri, Artikel, Tim).
- **Menghubungkan ke situs Next.js**: dashboard ini sudah siap dipakai
  sendiri dan punya API publik (`/api/services`, `/api/jobs`, `/api/articles`,
  `/api/facilities`, `/api/gallery`, `/api/team`, `/api/branches`,
  `/api/testimonials`, `/api/settings`) yang mengembalikan JSON. Untuk
  membuat situs Next.js benar-benar menarik data dari sini (bukan lagi dari
  file statis di `lib/*.ts`), langkah selanjutnya adalah mengganti setiap
  file tersebut dengan `fetch()` ke API ini saat build/request — ini sengaja
  belum dilakukan di sesi ini supaya situs yang sudah live tidak berisiko
  ikut berubah sebelum dashboard ini benar-benar dites dan diisi datanya.

## Deploy

Backend ini berdiri sendiri (binary Go + file SQLite + folder uploads),
bisa dijalankan di VPS mana pun (mis. Railway, Fly.io, VPS biasa) — beda
platform dari Vercel yang dipakai untuk situs Next.js sekarang. Untuk
build production:

```bash
go build -o gensakidz-backend .
./gensakidz-backend
```

Pastikan folder `uploads/` dan file `gensakidz.db` ikut di-backup/persist —
keduanya menyimpan semua data & foto yang diunggah lewat dashboard.
