import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";
import { cn } from "@/lib/utils";
import { Navbar } from "@/components/navbar";
import Providers from "@/components/providers";

const sfProDisplay = localFont({
  src: [
    { path: "../public/fonts/SFPRODISPLAYREGULAR.woff2", weight: "400", style: "normal" },
    { path: "../public/fonts/SFPRODISPLAYMEDIUM.woff2", weight: "500", style: "normal" },
    { path: "../public/fonts/SFPRODISPLAYBOLD.woff2", weight: "700", style: "normal" },
    { path: "../public/fonts/SFPRODISPLAYLIGHTITALIC.woff2", weight: "300", style: "italic" },
    { path: "../public/fonts/SFPRODISPLAYSEMIBOLDITALIC.woff2", weight: "600", style: "italic" },
    { path: "../public/fonts/SFPRODISPLAYTHINITALIC.woff2", weight: "100", style: "italic" },
    { path: "../public/fonts/SFPRODISPLAYULTRALIGHTITALIC.woff2", weight: "200", style: "italic" },
    { path: "../public/fonts/SFPRODISPLAYHEAVYITALIC.woff2", weight: "900", style: "italic" },
    { path: "../public/fonts/SFPRODISPLAYBLACKITALIC.woff2", weight: "900", style: "italic" },
  ],
  variable: "--font-sans",
});

export const metadata: Metadata = {
  title: "Rent",
  description: "",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={cn("h-full", "antialiased", sfProDisplay.variable, "font-sans")}
    >
      <body className="min-h-full flex flex-col">
        <Providers>
          {children}
          <Navbar />
        </Providers>
      </body>
    </html>
  );
}
