"use client";

import { useEffect } from "react";
import { AlertCircle, X } from "lucide-react";

import { cn } from "@/lib/utils";

// Notifikasi singkat untuk kegagalan aksi. Satu pesan, satu tempat, hilang
// otomatis — supaya error tidak menggeser layout kartu di bawahnya.
export function NoticeToast({
  open,
  title,
  message,
  onClose,
  duration = 6000,
}: {
  open: boolean;
  title: string;
  message: string;
  onClose: () => void;
  duration?: number;
}) {
  useEffect(() => {
    if (!open) return;
    const timer = setTimeout(onClose, duration);
    return () => clearTimeout(timer);
  }, [open, message, duration, onClose]);

  if (!open) return null;

  return (
    <div
      role="alert"
      aria-live="assertive"
      className={cn(
        "toast-in fixed bottom-6 right-6 z-50 flex w-[min(24rem,calc(100vw-3rem))] items-start gap-3",
        "rounded-xl border bg-popover p-4 text-popover-foreground shadow-lg",
      )}
    >
      <AlertCircle
        aria-hidden="true"
        className="mt-0.5 size-5 shrink-0 text-destructive"
      />
      <div className="min-w-0 flex-1 space-y-1">
        <p className="text-sm font-medium">{title}</p>
        <p className="text-sm break-words text-muted-foreground">{message}</p>
      </div>
      <button
        type="button"
        aria-label="Tutup notifikasi"
        onClick={onClose}
        className="-mr-1 -mt-1 rounded-md p-1 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
      >
        <X aria-hidden="true" className="size-4" />
      </button>
    </div>
  );
}
