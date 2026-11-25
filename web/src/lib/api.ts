import axios from "axios";

export const api = axios.create({
  // Same-origin: semua panggilan lewat rewrite Next (/api/* → API). Dengan
  // begitu cookie sesi tetap first-party walau web dibuka lewat
  // 127.0.0.1, IP LAN, atau port forwarding — bukan hanya localhost.
  baseURL: "",
  withCredentials: true, // kirim cookie access/refresh token
  headers: { "Content-Type": "application/json" },
});

// Satu refresh untuk semua request yang kena 401 bersamaan. Tanpa ini
// setiap request memicu panggilan /auth/refresh sendiri-sendiri.
let refreshPromise: Promise<void> | null = null;
// Setelah refresh gagal, sesi memang habis: jangan coba lagi sampai ada
// login/register/refresh yang sukses, supaya tidak terjadi badai 401/400.
let sessionExpired = false;

function refreshSession(): Promise<void> {
  if (!refreshPromise) {
    refreshPromise = api
      .post("/api/v1/auth/refresh")
      .then(() => {
        sessionExpired = false;
      })
      .catch((err) => {
        sessionExpired = true;
        throw err;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

api.interceptors.response.use(
  (res) => {
    const url = res.config?.url ?? "";
    if (
      url.includes("/auth/login") ||
      url.includes("/auth/register") ||
      url.includes("/auth/refresh")
    ) {
      sessionExpired = false;
    }
    return res;
  },
  async (error) => {
    const original = error.config;
    const url: string = original?.url ?? "";
    // Access token kedaluwarsa → coba refresh sekali (dibagi bersama
    // request lain), kecuali sesi sudah dinyatakan habis.
    if (
      error.response?.status === 401 &&
      original &&
      !original._retried &&
      !url.includes("/auth/refresh") &&
      !url.includes("/auth/login") &&
      !url.includes("/auth/logout") &&
      !sessionExpired
    ) {
      original._retried = true;
      try {
        await refreshSession();
        return api(original);
      } catch {
        // Refresh gagal → biarkan error 401 asli naik.
      }
    }
    return Promise.reject(error);
  },
);

// Ambil field `data` dari envelope backend; lempar error bila kosong
// supaya mutasi/query memasuki state error, bukan mengembalikan undefined.
export function unwrap<T>(envelope: { data?: { data?: T } }): T {
  const data = envelope.data?.data;
  if (data === undefined || data === null) {
    throw new Error("respons server tidak berisi data");
  }
  return data;
}

// Ambil pesan error dari envelope backend (axios error).
export function apiErrorMessage(err: unknown, fallback: string): string {
  const message = (
    err as { response?: { data?: { message?: unknown } } } | undefined
  )?.response?.data?.message;
  return typeof message === "string" && message ? message : fallback;
}
