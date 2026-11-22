# FIX_REPORT

Laporan perbaikan issue take-home assignment Meme Portal. Setiap issue punya
satu commit tersendiri dengan penjelasan akar masalah. Semua bukti di bawah
diambil dari aplikasi yang berjalan lokal (Postgres, Redis, Garage, API
`:8080`, web `:3000`).

- Repo: `meme-generator`
- Stack: Go + `net/http` + pgx + sqlc + Redis + Garage S3; Next.js App Router +
  TanStack Query + shadcn/ui
- Bukti "sebelum" diambil dari binary kode awal (`17c95de`) yang dijalankan di
  port `8081` dengan database yang sama.

## Ringkasan

| Issue | Akar masalah | Commit |
|---|---|---|
| TKT-101 | `ActiveByUser` tidak memeriksa `ends_at`, janitor expired hanya tiap jam | `821115d` |
| TKT-102 | `/subscription/status` dikirim `Cache-Control: public, max-age=86400` | `a736064` |
| TKT-103 | `clearAuthCookies` memakai `Path` yang salah; halaman login memakai `document.cookie` | `6319f5c` |
| TKT-104 | `staleTime: Infinity` pada query status | `8757cd0` |
| TKT-105 | Idempotency-Key dibuat ulang tiap klik; replay tidak dibatasi per user | `d9c1c21`, `701e887` |
| TKT-106 | Error konfirmasi pembayaran ditelan `catch` kosong | `6d9d20e` |
| TKT-107 | Query meme tidak memfilter `user_id` | `5350382` |
| TKT-108 | Tidak ada allowlist ekstensi/tipe/ukuran pada presigned upload | `ed76dcb` |
| TKT-109 | Cek kuota dan insert record tidak atomik | `60fd413` |

Bug tambahan yang ditemukan saat verifikasi:

| Bug | Commit |
|---|---|
| `GetByID` memindai 9 tujuan untuk 10 kolom → semua penyelesaian meme 500 | `a2d981f` |
| `Idempotency-Key` tidak diizinkan preflight CORS → tombol Berlangganan diam | `701e887` |

## TKT-101 — Masih bisa membuat meme padahal paket habis

**Gejala.** User membeli paket 1 hari, lewat masa aktif masih bisa upload.
Restart browser sempat normal lalu bisa lagi.

**Akar masalah.** `subscription.Repository.ActiveByUser` hanya memfilter
`status = 'active'` tanpa memeriksa `ends_at`, padahal komentar kodenya
mengklaim ada lazy check `ends_at > now()`. Janitor penanda expired berjalan
sekali per jam, sehingga selama jeda itu langganan kedaluwarsa tetap
"aktif" — baik untuk gate pembuatan meme maupun status di UI. Selain itu
kuota dihitung dengan jendela bergulir `now() - duration_days`, bukan sejak
langganan dimulai.

**Perubahan.**

- `api/internal/modules/subscription/repository.go`: `ActiveByUser` menambah
  `AND s.ends_at > now()`.
- `api/internal/modules/subscription/service.go`: `Status` menghitung
  `MemesUsed` sejak `subscriptions.starts_at`.

**Bukti sebelum** (`:8081`):

```
DB: status='active' | ends_at lampau: true
GET /api/v1/subscription/status
{"subscribed":true,"status":"active","planName":"Premium 1 Hari",...}
POST /api/v1/memes/upload
{"success":true,"code":30106,"message":"upload disiapkan",...}   [201]
```

**Bukti sesudah** (`:8080`, skenario sama):

```
GET /api/v1/subscription/status
{"subscribed":false,"memesUsed":0,"maxMemes":0}
POST /api/v1/memes/upload
{"success":false,"code":20002,"message":"butuh langganan aktif untuk membuat meme"}   [403]
```

## TKT-102 — Status langganan aneh, muncul di perangkat user lain

**Gejala.** Di WiFi kantor, layar rekan kerja menampilkan status dan kuota
milik user lain yang belum login.

**Akar masalah.** Endpoint per-user `/api/v1/subscription/status` dikirim
dengan `Cache-Control: public, max-age=86400`. Cache bersama (proxy/CDN)
menyimpan respons pertama dan menyajikannya ke pengguna lain selama sehari.

**Perubahan.** `subscription/mount.go`: header menjadi
`Cache-Control: private, no-store`.

**Bukti sebelum:**

```
$ curl -D - -o /dev/null -b cookie.txt localhost:8081/api/v1/subscription/status
Cache-Control: public, max-age=86400
```

**Bukti sesudah:**

```
$ curl -D - -o /dev/null -b cookie.txt localhost:8080/api/v1/subscription/status
Cache-Control: private, no-store
```

## TKT-103 — Logout tidak berfungsi dan halaman muter-muter

**Akar masalah.** Dua bug yang saling melengkapi:

1. `auth.Handler.clearAuthCookies` menulis cookie penghapus dengan
   `Path=/api/v1/auth`, sedangkan cookie sesi dibuat dengan `Path=/`. Browser
   menganggap keduanya cookie berbeda, jadi cookie sesi tidak terhapus.
2. Halaman `/login` memakai `document.cookie.includes("meme_access_token")`
   untuk redirect, sehingga cookie sisa membuat `/login` melempar ke `/`,
   sementara `/` melempar balik ke `/login` karena sesi sudah direvoke.

**Perubahan.**

- `auth/handler.go`: hapus cookie lewat `cookie.Clear` (Path `/`).
- `web/src/app/login/page.tsx`: redirect berdasarkan `useMe()` (`/auth/me`),
  bukan `document.cookie`.

**Bukti sesudah:**

```
POST /api/v1/auth/logout
Set-Cookie: meme_access_token=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax
Set-Cookie: meme_refresh_token=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax
GET /api/v1/auth/me → 401
isi cookie jar setelah logout: 0
```

UI (Chromium): klik **Keluar** → halaman berhenti di `/login`, navigasi ulang
ke `/login` tetap di `/login` (tidak ada ping-pong).

## TKT-104 — Navbar bilang belum langganan sampai hard refresh

**Akar masalah.** `useSubscriptionStatus` memakai
`staleTime: Number.POSITIVE_INFINITY`, jadi setelah pembayaran query tidak
pernah diambil ulang. Diperparah header cache publik pada TKT-102.

**Perubahan.** `web/src/lib/queries.ts`: `staleTime` menjadi 10 detik; query
sudah di-invalidate pada `useCheckout` dan `useConfirmPayment`.

**Bukti sesudah** (UI, tanpa hard refresh): setelah klik **Bayar Sekarang**,
halaman Paket langsung menampilkan

```
Status: Premium 1 Hari (aktif sampai 2026-11-23T15:22:56Z)
Kuota meme: 0 / 10 (sisa 10)
```

## TKT-105 — Tagihan dobel, checkout ganda

**Akar masalah.** `useCheckout` membuat `crypto.randomUUID()` **di dalam**
`mutationFn`, jadi setiap klik (termasuk double-click) dianggap maksud
berbeda. Backend pun melakukan replay lookup idempotency tanpa filter user,
dan dua request dengan key sama yang tiba bersamaan menghasilkan 500
(unique index). Ditemukan juga: header `Idempotency-Key` tidak diizinkan
preflight CORS sehingga checkout dari browser gagal total.

**Perubahan.**

- `web/src/lib/queries.ts`: satu key per paket disimpan di `useRef` dan
  dipakai ulang untuk retry; tombol dinonaktifkan selama pending.
- `subscription/repository.go` + `service.go`: replay lookup dibatasi
  `s.user_id = $2`; pelanggaran unique index (23505) mengembalikan checkout
  yang sudah ada, bukan 500.
- `middleware/cors.go`: `Idempotency-Key` masuk `Access-Control-Allow-Headers`.

**Bukti sebelum** (dua klik, key baru per klik):

```
jumlah tagihan untuk 2 klik: 2
```

**Bukti sesudah** (dua request paralel, key sama):

```
respons 1 & 2 → subscription id identik, HTTP 201
jumlah tagihan: 1
```

**Bukti UI:** double-click **Berlangganan** → satu baris pembayaran
(`count(*)=1`), sekali redirect ke halaman checkout.

## TKT-106 — Klik bayar tidak terjadi apa-apa

**Akar masalah.** `onPay` menelan semua error (`catch {}`) sehingga gagal
bayar tidak menampilkan apa pun, dan tombol tidak berubah saat request
berjalan.

**Perubahan.** `web/src/app/checkout/[reference]/page.tsx`: error ditampilkan
lewat `describeError`, tombol menampilkan `Memproses…` dan dinonaktifkan
selama request serta setelah sukses.

**Bukti sesudah** (UI, membayar referensi yang sudah lunas):

```
Pembayaran gagal: HTTP 402: pembayaran gagal
```

## TKT-107 — Bisa melihat dan menghapus meme user lain

**Akar masalah.** `GetByID`, `SetResultURL`, dan `Delete` hanya memfilter
kolom `id`, bukan `user_id`. Siapa pun yang mengganti id di URL bisa
menyelesaikan atau menghapus meme user lain.

**Perubahan.** `meme/repository.go`: semua query tersebut menambah
`AND user_id = $N`; id milik user lain menghasilkan 404 (`ErrNotFound`).

**Bukti sebelum** (user B memakai id milik user A):

```
DELETE /api/v1/memes/<id-A> → {"success":true,"code":30105,"message":"meme dihapus"} [200]
POST /api/v1/memes/<id-A>/complete → {"success":true,...,"title":"race"} [200]
```

**Bukti sesudah:**

```
DELETE      → {"code":30001,"message":"meme tidak ditemukan"} [404]
complete    → {"code":30001,"message":"meme tidak ditemukan"} [404]
result-upload → {"code":30001,"message":"meme tidak ditemukan"} [404]
pemilik result-upload → 200
```

## TKT-108 — Ada installer GTA V di bucket S3

**Akar masalah.** `PrepareUpload` tidak memakai helper `allowedExt` yang sudah
ada, memakai `Content-Type` kiriman klien apa adanya, dan melakukan presign
tanpa batas ukuran. Siapa pun bisa mengunggah `.zip` 90 GB, `.exe`, atau
`.html` ke bucket publik.

**Perubahan.**

- `meme/validate.go` (baru): `validateBaseImage` menolak ekstensi di luar
  `jpg/jpeg/png/gif/webp` (400) dan ukuran > 8 MiB (413); `allowedType`
  memetakan ekstensi ke Content-Type kanonik.
- `meme/handler.go`: validasi dijalankan sebelum presign; `Content-Type`
  dari klien tidak dipakai lagi.
- `meme/service.go`: `ContentType` dan `ContentLength` ikut ditandatangani
  pada presigned PUT, sehingga upload dengan tipe/ukuran berbeda ditolak S3.

**Bukti sebelum:**

```
POST /api/v1/memes/upload {"filename":"GTA.V.Full.Installer.zip","size":90000000000}
{"success":true,...,"uploadUrl":".../base/e38107b2521b07f15e726e2e.zip?...&X-Amz-..."} [201]
```

**Bukti sesudah:**

```
filename=a.zip,size=100        → {"code":30002,"message":"file gambar tidak valid"} [400]
filename=evil.html             → [400]
filename=malware.exe           → [400]
filename=a.png,size=90000000000 → {"code":30003,"message":"ukuran gambar melebihi batas"} [413]
filename=a.png,size=100        → 201, presigned URL
PUT 100 byte, Content-Type image/png  → 200
PUT 200 byte (ukuran beda)            → 403
PUT 100 byte, Content-Type text/html  → 403
```

## TKT-109 — Kuota 20 meme tapi bisa membuat 25

**Akar masalah.** Pengecekan kuota (`MemesUsed >= MaxMemes`) dan `INSERT`
record meme berjalan sebagai dua query terpisah tanpa transaksi. Enam tab
mengirim request bersamaan, semuanya membaca hitungan yang sama sebelum ada
yang commit, lalu semuanya lolos.

**Perubahan.**

- `meme/repository.go`: `CreateWithQuota` menjalankan hitung kuota dan
  insert dalam satu transaksi dengan `pg_advisory_xact_lock(hashtext(user_id))`
  sehingga request paralel user yang sama diproses berurutan.
- `meme/service.go` + `subscription/service.go`: batas kuota dan awal periode
  (`starts_at`) diambil dari langganan aktif lalu diserahkan ke transaksi itu;
  `repository.Create` lama dihapus.

**Bukti sebelum** (20 request paralel, paket `premium_1d` kuota 10):

```
15 201
 5 403
memesUsed=16 maxMemes=10
```

**Bukti sesudah** (20 request paralel, kuota 10):

```
10 201
10 403
memesUsed=10 maxMemes=10
```

## Bug tambahan saat verifikasi

1. **`GetByID` gagal scan** (`a2d981f`). Query memilih 10 kolom
   (`... created_at, updated_at`) tetapi `Scan` hanya menyiapkan 9 tujuan,
   sehingga `result-upload` dan `complete` selalu 500:

   ```
   level=ERROR msg="meme handler: internal error"
   err="number of field descriptions must equal number of destinations, got 10 and 9"
   ```

   Setelah perbaikan, alur lengkap berjalan: `PUT base=200`, `PUT result=200`,
   `complete → 200 {"message":"meme selesai"}`.

2. **Preflight CORS menolak `Idempotency-Key`** (`701e887`) — lihat TKT-105.

## Verifikasi

- Setiap commit diuji dapat di-build (`go build ./...`) dari arsip commit-nya
  masing-masing.
- `go vet ./...` bersih, `tsc --noEmit` bersih, `next lint` tanpa warning.
- Alur utama diuji pada aplikasi berjalan: register → checkout → bayar →
  status aktif → upload gambar dasar → upload hasil → complete → daftar meme →
  logout.

## Catatan

- Payment gateway tetap mock (sesuai README); tidak diubah ke provider nyata.
- Batas ukuran gambar dasar 8 MiB (konstanta `maxBaseImageBytes`).
- File lama yang sudah terlanjur ada di bucket (mis. `GTA.V.Full.Installer.zip`)
  tidak dihapus otomatis oleh perbaikan ini. Pembersihan bucket dan pemulihan
  meme yang terhapus karena TKT-107 perlu dilakukan manual di storage/DB.
- Pemulihan data: satu meme user lain memang terhapus pada bukti reproduksi.
  Karena `DELETE` bersifat hard delete, data tersebut tidak bisa dipulihkan
  dari aplikasi; restore harus dari backup database.
