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

## Menjalankan via Docker

Cara paling praktis untuk menjalankan (lokal maupun deploy) — sudah diuji
build, restart, dan redeploy penuh (`down` lalu `up`), datanya tetap aman
tersimpan di volume Docker:

```bash
cd backend
cp .env.example .env    # lalu edit .env — minimal ganti ADMIN_PASSWORD
docker compose up -d --build
```

Server jalan di `http://localhost:8097` (port host bisa diganti lewat
`HOST_PORT` di `.env`, dipilih bukan 8080 dari awal supaya tidak gampang
bentrok dengan service lain di komputer/server yang sama). Data (database
SQLite + foto yang diunggah) tersimpan di Docker volume `gensakidz_data`
— aman biarpun container dihapus dan dibuat ulang (mis. setiap kali
redeploy).

`.env` (lihat `.env.example` untuk daftar lengkap variabelnya) tidak
ikut ke git — setiap deploy baru wajib bikin `.env` sendiri dan ganti
`ADMIN_PASSWORD` dari nilai default.

Perintah lain yang berguna:

```bash
docker compose logs -f backend   # lihat log
docker compose restart backend   # restart tanpa kehilangan data
docker compose down              # hentikan & hapus container (volume tetap ada)
```

`Dockerfile` pakai multi-stage build (`CGO_ENABLED=0`, image akhir
`alpine`, jalan sebagai user non-root) — image jadi kecil dan tidak
butuh Go/toolchain apa pun di server tujuan, cukup Docker saja.

## Deploy ke hosting gratis (Fly.io)

Backend ini butuh **disk permanen** (untuk file SQLite + foto upload),
jadi platform serverless murni (Vercel, Netlify, tier gratis
Render/Railway) tidak cocok — datanya akan hilang tiap restart. Fly.io
mendukung volume permanen dan punya jatah gratis yang cukup untuk skala
klinik kecil begini:

```bash
brew install flyctl        # atau lihat fly.io/docs/getting-started
fly auth login
cd backend
fly launch --no-deploy      # pilih region Singapore (sin) biar dekat
fly volumes create gensakidz_data --size 1 --region sin
fly deploy
fly secrets set ADMIN_EMAIL=owner@gensakidz.com ADMIN_PASSWORD="kata-sandi-kuat"
```

Saat `fly launch` menanyakan konfigurasi, pastikan `fly.toml` yang
dihasilkan me-mount volume ke `/app/data` (samakan dengan `DB_PATH` dan
`UPLOADS_DIR` di atas) — detail lengkap ada di dokumentasi Fly.io karena
`fly launch` bersifat interaktif dan hasilnya tergantung pilihan saat
setup.
