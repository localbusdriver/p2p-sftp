import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const InterFont = Inter({
  weight: ["400", "500", "600", "700", "800", "900"],
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Create Next App",
  description: "Generated by create next app",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${InterFont.className} antialiased`}>{children}</body>
    </html>
  );
}
