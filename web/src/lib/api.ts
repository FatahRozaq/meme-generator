import axios from "axios";

export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  withCredentials: true, // kirim cookie access/refresh token
  headers: { "Content-Type": "application/json" },
});

api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config;
    // Access token kedaluwarsa → coba refresh sekali lewat proxy rewrite.
    if (
      error.response?.status === 401 &&
      !original._retried &&
      !original.url?.includes("/auth/refresh")
    ) {
      original._retried = true;
      try {
        await api.post("/api/v1/auth/refresh");
        return api(original);
      } catch {
        // refresh gagal → biarkan error 401 asli naik.
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
