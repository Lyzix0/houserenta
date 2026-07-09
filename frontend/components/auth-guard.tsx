"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getMe } from "@/lib/auth";

type AuthGuardProps = {
  mode: "guest" | "auth";
  children: React.ReactNode;
};

export function AuthGuard({ mode, children }: AuthGuardProps) {
  const router = useRouter();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    getMe()
      .then(() => {
        if (mode === "guest") {
          router.replace("/home");
        } else {
          setChecking(false);
        }
      })
      .catch(() => {
        if (mode === "auth") {
          router.replace("/auth/login");
        } else {
          setChecking(false);
        }
      });
  }, [mode, router]);

  if (checking) {
    return (
      <div className="flex justify-center items-center h-40">
        <p className="text-sm text-muted-foreground font-medium">Загрузка...</p>
      </div>
    );
  }

  return <>{children}</>;
}
