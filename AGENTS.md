# AGENTS.md — Meme Portal

Panduan alur kerja dan standarisasi untuk agent/developer yang mengerjakan project ini.
Baca README.md untuk overview stack.

## Project

Meme generator berbasis langganan. User subscribe → upload gambar → tambah teks (top/bottom)
→ browser render via canvas → upload hasil ke S3.

- `api/` — Go (stdlib `net/http`, pgx/v5, sqlc, golang-migrate, JWT, Redis, Garage S3)
- `web/` — Next.js App Router (TanStack Query/Table/Form, shadcn/ui, axios, Bun)
- Infra lokal: `compose.yaml` (postgres 17, redis 7, garage v2.4.1) + `process-compose.yaml`

## Skills yang wajib dipakai

| Skill | Kapan |
|---|---|
| `/skill:ponytail` | Semua coding task: perubahan minimal, YAGNI, jangan tambah abstraksi |
| `/skill:impeccable` | Sentuh UI frontend (web/) — polish, konsistensi design system |
| `/skill:antislop` | Semua kode yang dihasilkan — hindari pola/slop khas AI |
| `/skill:antislop-code` | Bersihkan komentar/junk kode AI di file yang diedit |
| `/skill:antislop-ui` | Komponen UI: warna, layout, state, aksesibilitas, motion |

## Alur pengerjaan

### 0. Persiapan (sekali)

```sh
cp .env.example .env
make setup   # cek & pasang docker, go, bun, migrate, sqlc, air, bun install
```

Windows (host ini): `make`, `bun`, `process-compose` belum terpasang. Install manual:
- make: via scoop (`scoop install make`) atau gunakan perintah Makefile langsung
- bun: `powershell -c "irm bun.sh/install.ps1 | iex"`
- process-compose: download binary dari https://github.com/F1bonacc1/process-compose/releases

### 1. Menjalankan dev

```sh
make dev   # infra + migrate + api :8080 + web :3000 via process-compose (TUI di :8099)
```

Catatan Windows:
- Binary air dikonfigurasi sebagai `tmp/api.exe` di `api/.air.toml` (Windows menolak
  mengeksekusi binary tanpa ekstensi `.exe`).
- Kalau muncul log `Path not found for process compose config home`, buat folder
  `%LOCALAPPDATA%\process-compose` (sekali saja).

Tanpa make, jalankan manual urut:
1. `docker compose up -d --wait postgres redis garage`
2. `MSYS_NO_PATHCONV=1 docker compose exec garage /garage bucket website --allow meme-portal`
   (agar URL publik gambar bisa diakses browser; prefix `MSYS_NO_PATHCONV=1`
   hanya perlu di Git Bash Windows agar `/garage` tidak diubah jadi path Windows)
3. `cd api && migrate -path migrations -database "$DATABASE_URL" up`
4. `cd api && air` (hot reload, `api/.air.toml`)
5. `cd web && bun run dev`

### 2. Alur kerja backend

Modul pattern (wajib ikuti, lihat `api/internal/modules/auth/` sebagai contoh):
`handler.go` (HTTP), `service.go` (bisnis), `repository.go` (DB via sqlc),
`model.go`, `dto.go`, `errors.go`, `codes.go`, `mount.go` (rute).

1. Ubah schema → `make migrate-new name=x` → tulis SQL di `api/migrations/` (up+down)
2. Tulis query di `api/internal/db/queries/*.sql` → `make sqlc-generate`
3. Tambah/ubah handler+service+repository sesuai pattern modul
4. Response selalu pakai envelope `api/internal/response/response.go`
5. JWT auth: `internal/middleware/auth.go`; refresh token disimpan di Redis (revokeable)

### 3. Alur kerja frontend

1. Tipe API baru → `web/src/lib/types.ts`
2. Hook data → `web/src/lib/queries/` (TanStack Query), client axios di `web/src/lib/api.ts`
3. Halaman → `web/src/app/<route>/page.tsx` (App Router)
4. Komponen UI → pakai shadcn (`web/src/components/ui/`), jangan styling ad-hoc
5. State form → TanStack Form

### 4. Verifikasi sebelum selesai

```sh
make lint   # lint api + web
make test   # test api
make check  # lint + build
```

Jalankan surface yang diubah (API endpoint via curl / halaman web via browser), bukan cuma build.

## Konvensi

- Backend: env-based config (`api/internal/config/config.go`), jangan hardcode
- Semua env baru: tambahkan ke `.env.example` dan `process-compose.yaml` bila perlu
- S3: bucket lokal `meme-portal` via Garage; public URL `S3_PUBLIC_BASE_URL`
- Payment gateway mock di modul subscription — jangan ganti ke provider nyata tanpa diminta
- Migrasi: seq, selalu ada `.up.sql` + `.down.sql`, jangan edit migrasi yang sudah di-apply
- Bahasa komentar/dokumen project: Indonesia

## Anti-pattern (jangan)

- Jangan tambah dependency baru tanpa perlu (ponytail)
- Jangan bikin abstraksi/layer baru; ikuti pattern modul yang ada
- Jangan commit `.env`, `tsconfig.tsbuildinfo`
- Jangan skip `.down.sql`
