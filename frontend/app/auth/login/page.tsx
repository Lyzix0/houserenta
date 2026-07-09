import { Suspense } from "react";
import { LoginForm } from "@/components/login-form";
import { AuthGuard } from "@/components/auth-guard";

export default function Page() {
  return (
    <AuthGuard mode="guest">
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
        <div className="w-full max-w-sm">
          <Suspense fallback={
            <div className="flex justify-center items-center h-40">
              <p className="text-sm text-muted-foreground font-medium">Загрузка...</p>
            </div>
          }>
            <LoginForm />
          </Suspense>
        </div>
      </div>
    </AuthGuard>
  );
}
