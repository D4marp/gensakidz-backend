# GenSA Kidz — Backend & Dashboard

Backend Go sederhana + dashboard admin untuk mengelola semua konten yang saat
ini tampil statis di gensakidz.com (layanan, karir, artikel, fasilitas,
galeri, tim, cabang, testimoni, dan teks umum seperti visi/misi). Sekali
dijalankan, seluruh konten yang sudah ada di situs otomatis dimasukkan ke
database — dashboard tidak dimulai dari kosong.

Database: **MySQL**. Foto yang diunggah lewat dashboard disimpan sebagai
file biasa di folder `uploads/`.

## Menjalankan via Docker (disarankan)

Cara paling praktis — sudah diuji build, restart, dan redeploy penuh
(`down` lalu `up`), datanya tetap aman tersimpan di volume Docker, termasuk
database MySQL-nya:

```bash
cd backend
cp .env.example .env    # lalu edit .env — minimal ganti semua password
docker compose up -d --build
```

Ini menjalankan **dua container**: `mysql` (database) dan `backend` (Go
app) — backend otomatis menunggu MySQL siap sebelum jalan.

Server jalan di `http://localhost:8097` (port host bisa diganti lewat
`HOST_PORT` di `.env`, dipilih bukan 8080 dari awal supaya tidak gampang
bentrok dengan service lain di komputer/server yang sama).

`.env` (lihat `.env.example` untuk daftar lengkap variabelnya) tidak ikut
ke git — setiap deploy baru wajib bikin `.env` sendiri dan ganti semua
password dari nilai default.

Perintah lain yang berguna:

```bash
docker compose logs -f backend   # lihat log backend
docker compose logs -f mysql     # lihat log database
docker compose restart backend   # restart tanpa kehilangan data
docker compose down              # hentikan & hapus container (volume tetap ada)
```

`Dockerfile` pakai multi-stage build (`CGO_ENABLED=0`, image akhir
`alpine`, jalan sebagai user non-root) — image jadi kecil dan tidak
butuh Go/toolchain apa pun di server tujuan, cukup Docker saja.

## Menjalankan tanpa Docker

Butuh MySQL server sendiri yang sudah jalan (lokal atau eksternal):

```bash
cd backend
export DB_HOST=127.0.0.1 DB_PORT=3306 DB_USER=gensakidz DB_PASSWORD=... DB_NAME=gensakidz
go run .
```

Server jalan di `http://localhost:8080`, otomatis redirect ke `/admin`.

Login default (**wajib diganti** — lihat bagian Konfigurasi):
- Email: `admin@gensakidz.com`
- Kata sandi: `gensakidz2026`

## Struktur

- `main.go` — routing & start server
- `db.go` — koneksi MySQL & skema (dibuat otomatis lewat `CREATE TABLE IF NOT EXISTS`)
- `models.go` — struct data (Service, Job, Article, Facility, dst.)
- `repo.go` — query database per jenis konten
- `seed.go` — data awal (migrasi dari konten statis situs saat ini)
- `handlers.go` — halaman dashboard admin (list, tambah, ubah, hapus)
- `api.go` — API publik JSON read-only (`/api/*`) untuk dikonsumsi frontend
- `auth.go` — login sederhana (cookie session + bcrypt)
- `upload.go` — upload foto, disimpan di folder `uploads/`
- `templates/`, `static/` — tampilan dashboard (HTML + CSS polos, tanpa build step)

## Konfigurasi (environment variable)

| Variabel | Default | Keterangan |
|---|---|---|
| `PORT` | `8080` | Port server (di dalam container) |
| `HOST_PORT` | `8097` | Port di komputer/server host (khusus `docker-compose.yml`) |
| `DB_HOST` | `127.0.0.1` | Host MySQL (di Docker: otomatis `mysql`, nama service) |
| `DB_PORT` | `3306` | Port MySQL |
| `DB_USER` | `gensakidz` | User MySQL |
| `DB_PASSWORD` | `gensakidz` | Password MySQL — **wajib diganti** |
| `DB_NAME` | `gensakidz` | Nama database |
| `ADMIN_EMAIL` | `admin@gensakidz.com` | Email login admin (dipakai hanya saat pertama kali dijalankan) |
| `ADMIN_PASSWORD` | `gensakidz2026` | Password admin (dipakai hanya saat pertama kali dijalankan) — **wajib diganti** |
| `UPLOADS_DIR` | `uploads` | Folder penyimpanan foto yang diunggah |

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
  file tersebut dengan `fetch()` ke API ini saat build/request.

## Deploy ke hosting gratis (Fly.io)

Backend + MySQL ini butuh **disk permanen**, jadi platform serverless murni
(Vercel, Netlify, tier gratis Render/Railway) tidak cocok — datanya akan
hilang tiap restart. Fly.io mendukung volume permanen dan punya jatah
gratis yang cukup untuk skala klinik kecil begini. Karena setup-nya ada
dua container (backend + MySQL), cara termudah adalah pakai
[Fly's Docker Compose support](https://fly.io/docs/launch/compose/) atau
deploy MySQL terpisah sebagai [managed MySQL add-on](https://fly.io/docs/launch/db/)
lalu arahkan `DB_HOST` dkk ke sana:

```bash
brew install flyctl        # atau lihat fly.io/docs/getting-started
fly auth login
cd backend
fly launch --no-deploy      # pilih region Singapore (sin) biar dekat
fly volumes create gensakidz_uploads --size 1 --region sin
fly deploy
fly secrets set ADMIN_EMAIL=owner@gensakidz.com ADMIN_PASSWORD="kata-sandi-kuat" \
  DB_HOST=... DB_USER=... DB_PASSWORD=... DB_NAME=...
```

Detail lengkap ada di dokumentasi Fly.io karena `fly launch` bersifat
interaktif dan hasilnya tergantung pilihan saat setup.
