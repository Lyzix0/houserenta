"use client";

import { useEffect, useState, useRef, use } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/use-user";
import { applyForProperty } from "@/lib/properties";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function ApplyPage({ params }: { params: Promise<{ propertyId: string }> }) {
  const router = useRouter();
  const { data: user, isLoading } = useUser();
  const { propertyId } = use(params);

  const [status, setStatus] = useState<"idle" | "applying" | "success" | "error">("idle");
  const [errorMessage, setErrorMessage] = useState("");
  const appliedRef = useRef(false);

  useEffect(() => {
    if (isLoading) return;

    if (!user) {
      router.push(`/auth/signup?role=tenant&next=/apply/${propertyId}`);
      return;
    }

    const autoApply = async () => {
      if (appliedRef.current) return;
      appliedRef.current = true;
      setStatus("applying");

      try {
        await applyForProperty(propertyId);
        setStatus("success");
        setTimeout(() => {
          router.push("/home");
        }, 3000);
      } catch (err: any) {
        if (err.message && err.message.includes("уже откликнулись")) {
          setStatus("success");
          setTimeout(() => {
            router.push("/home");
          }, 3000);
        } else {
          setStatus("error");
          setErrorMessage(err.message || "Ошибка при отправке отклика");
        }
      }
    };

    autoApply();
  }, [user, isLoading, propertyId, router]);

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10 bg-slate-50 dark:bg-slate-900">
      <div className="w-full max-w-md">
        <Card className="border-slate-200 shadow-lg">
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Оформление отклика на аренду</CardTitle>
            <CardDescription>Заявка на заселение по приглашению владельца</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col items-center justify-center py-6 text-center text-xs md:text-sm space-y-4">
            {isLoading || status === "idle" || status === "applying" ? (
              <div className="flex flex-col items-center gap-3">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
                <p className="text-muted-foreground font-semibold">Инициализация и отправка заявки...</p>
              </div>
            ) : status === "success" ? (
              <div className="space-y-4">
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 text-emerald-600">
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <h3 className="text-lg font-bold text-slate-800 dark:text-white">Отклик успешно отправлен!</h3>
                <p className="text-muted-foreground">Вы успешно откликнулись на это предложение. Арендодатель сможет заселить вас в один клик.</p>
                <p className="text-[10px] text-muted-foreground animate-pulse">Перенаправление в личный кабинет через 3 сек...</p>
                <Button onClick={() => router.push("/home")} className="w-full bg-slate-900 text-white font-bold rounded-xl mt-2">
                  Перейти в личный кабинет
                </Button>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-rose-100 text-rose-600">
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
                <h3 className="text-lg font-bold text-slate-800 dark:text-white">Не удалось отправить заявку</h3>
                <p className="text-rose-600 font-semibold">{errorMessage}</p>
                <p className="text-muted-foreground text-xs">Убедитесь, что ссылка правильная, или перейдите в личный кабинет.</p>
                <Button onClick={() => router.push("/home")} className="w-full bg-slate-900 text-white font-bold rounded-xl mt-2">
                  На главную
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
