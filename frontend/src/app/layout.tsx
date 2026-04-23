import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { ToastProvider } from "@/lib/toast/ToastContext";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "HANetwork · Billing",
  description: "VPS/Proxy billing platform",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="h-full">
      <body className={`${inter.className} h-full overflow-hidden bg-[#F5F6F7] font-medium text-base`}>
        <ToastProvider>{children}</ToastProvider>
      </body>
    </html>
  );
}
