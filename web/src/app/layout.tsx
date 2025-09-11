import type { Metadata } from "next";
import "./globals.css";
import Script from "next/script";
import { AuthProvider } from "@/contexts/auth-context";
export const metadata: Metadata = {
  title: "Mindful — Track, Care, Grow",
  description: "Track your mental health, mood, and daily habits with Mindful"
};
export default function RootLayout({
  children
}: Readonly<{children: React.ReactNode;}>) {
  return (
    <html lang="en" data-oid="b1jifow" suppressHydrationWarning>
      <body className="antialiased" data-oid="hqec2:a">
        <AuthProvider data-oid="e2cdiiq">{children}</AuthProvider>
        <Script
          type="module"
          strategy="afterInteractive"
          src="https://cdn.jsdelivr.net/gh/onlook-dev/onlook@main/apps/web/client/public/onlook-preload-script.js"
          data-oid="zxbuq9a" />

      </body>
    </html>);

}