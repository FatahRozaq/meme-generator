"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import { useState } from "react";

// Error 4xx tidak akan berubah bila dicoba ulang (khususnya 401 setelah
// sesi habis); hanya gangguan server yang layak diulang sekali.
function shouldRetry(failureCount: number, error: unknown) {
  if (isAxiosError(error)) {
    const status = error.response?.status ?? 0;
    if (status >= 400 && status < 500) return false;
  }
  return failureCount < 1;
}

export function Providers({ children }: { children: React.ReactNode }) {
  const [client] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 30_000,
            refetchOnWindowFocus: false,
            retry: shouldRetry,
          },
        },
      }),
  );

  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}
