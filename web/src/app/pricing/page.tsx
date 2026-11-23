"use client";

// Pricing: daftar paket dari backend; checkout lewat payment gateway mock.

import { useCallback, useState } from "react";
import { useRouter } from "next/navigation";
import { Check } from "lucide-react";
import { format, parseISO } from "date-fns";
import { id as localeId } from "date-fns/locale";

import { useCheckout, usePlans, useSubscriptionStatus } from "@/lib/queries";
import { apiMessage, describeError } from "@/lib/errors";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { NoticeToast } from "@/components/notice-toast";

function formatEndsAt(value?: string) {
  if (!value) return "";
  const date = parseISO(value);
  if (Number.isNaN(date.getTime())) return value;
  return format(date, "d MMM yyyy, HH.mm", { locale: localeId });
}

function formatIdr(n: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(n);
}

export default function PricingPage() {
  const plans = usePlans();
  const status = useSubscriptionStatus();
  const checkout = useCheckout();
  const router = useRouter();

  const [notice, setNotice] = useState<string | null>(null);
  const resetCheckout = checkout.reset;
  const closeNotice = useCallback(() => {
    setNotice(null);
    resetCheckout();
  }, [resetCheckout]);

  const activeCode = status.data?.subscribed
    ? (status.data.planCode ?? null)
    : null;

  return (
    <main className="mx-auto max-w-4xl space-y-8 p-8">
      <h1 className="text-3xl font-bold">Paket Berlangganan</h1>
      {status.data && (
        <div className="space-y-1 text-sm text-muted-foreground">
          <p>
            Status:{" "}
            {status.data.subscribed
              ? `${status.data.planName} (aktif sampai ${formatEndsAt(status.data.endsAt)})`
              : "belum berlangganan"}
          </p>
          {status.data.subscribed && (
            <p>
              Kuota meme: {status.data.memesUsed}
              {status.data.maxMemes > 0
                ? ` / ${status.data.maxMemes} (sisa ${status.data.maxMemes - status.data.memesUsed})`
                : " / unlimited"}
            </p>
          )}
        </div>
      )}
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {plans.data?.map((plan) => {
          const isActive = plan.code === activeCode;
          return (
            <Card
              key={plan.id}
              aria-current={isActive ? "true" : undefined}
              className={cn(isActive && "border-foreground/25 shadow-sm")}
            >
              <CardHeader>
                <CardTitle>{plan.name}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <p className="text-2xl font-bold">{formatIdr(plan.priceIdr)}</p>
                <p className="text-sm text-muted-foreground">
                  {plan.durationDays} hari ·{" "}
                  {plan.maxMemes === 0
                    ? "unlimited meme"
                    : `${plan.maxMemes} meme`}
                </p>
                {isActive ? (
                  <Button className="w-full" variant="secondary" disabled>
                    <Check aria-hidden="true" />
                    Paket Aktif
                  </Button>
                ) : (
                  <Button
                    className="w-full"
                    disabled={checkout.isPending}
                    onClick={() =>
                      checkout.mutate(plan.code, {
                        onSuccess: (res) =>
                          router.push(`/checkout/${res.payment.reference}`),
                        onError: (err) => setNotice(apiMessage(err, describeError(err))),
                      })
                    }
                  >
                    {checkout.isPending ? "Memproses…" : "Berlangganan"}
                  </Button>
                )}
              </CardContent>
            </Card>
          );
        })}
      </div>

      <NoticeToast
        open={notice !== null}
        title="Tidak bisa berlangganan"
        message={notice ?? ""}
        onClose={closeNotice}
      />
    </main>
  );
}
