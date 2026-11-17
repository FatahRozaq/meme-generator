import type { Metadata } from "next";
import { Anton } from "next/font/google";
import { Providers } from "@/components/providers";
import { SiteHeader } from "@/components/site-header";
import "./globals.css";

// Anton: fallback gratis untuk Impact (bold, condensed) agar tetap tampil
// konsisten di OS yang tidak punya font Impact (mis. Linux).
const anton = Anton({
  weight: "400",
  subsets: ["latin"],
  variable: "--font-anton",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Meme Portal",
  description: "Meme generator berlangganan dengan font Impact",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="id" className={anton.variable}>
      <body className="antialiased">
        <Providers>
          <SiteHeader />
          {children}
        </Providers>
      </body>
    </html>
  );
}
