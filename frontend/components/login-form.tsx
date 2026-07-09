"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { HugeiconsIcon } from "@hugeicons/react";
import { ViewIcon, ViewOffIcon } from "@hugeicons/core-free-icons";

import { cn } from "@/lib/utils";
import { login } from "@/lib/auth";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { PhoneInput } from "@/components/ui/phone-input";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [email, setEmail] = useState("");
  const [loginMethod, setLoginMethod] = useState<"email" | "phone">("email");
  const [password, setPassword] = useState("");

  const [showPassword, setShowPassword] = useState(false);

  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  const router = useRouter();
  const searchParams = useSearchParams();

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setErrorMessage("");
    setSuccessMessage("");

    try {
      await login({ email, password });

      setSuccessMessage("Вход выполнен успешно.");

      const nextPath = searchParams.get("next") || "/home";
      router.push(nextPath);
    } catch (err) {
      setErrorMessage(
        err instanceof Error
          ? err.message
          : "Ошибка соединения. Попробуйте позже."
      );
    }
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Авторизация</CardTitle>
          <CardDescription>
            Войдите в аккаунт, чтобы продолжить
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleLogin}>
            <FieldGroup>
              {loginMethod === "email" ? (
                <Field>
                  <div className="flex justify-between items-center">
                    <FieldLabel htmlFor="email">Почта</FieldLabel>
                    <button
                      type="button"
                      onClick={() => {
                        setLoginMethod("phone");
                        setEmail("");
                      }}
                      className="text-xs text-indigo-600 hover:underline font-semibold"
                    >
                      Войти по телефону
                    </button>
                  </div>

                  <Input
                    id="email"
                    type="email"
                    placeholder="example@mail.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                  />
                </Field>
              ) : (
                <Field>
                  <div className="flex justify-between items-center">
                    <FieldLabel htmlFor="phone">Телефон</FieldLabel>
                    <button
                      type="button"
                      onClick={() => {
                        setLoginMethod("email");
                        setEmail("");
                      }}
                      className="text-xs text-indigo-600 hover:underline font-semibold"
                    >
                      Войти по почте
                    </button>
                  </div>

                  <PhoneInput
                    id="phone"
                    placeholder="Введите номер телефона"
                    value={email}
                    onChange={(val) => setEmail(val || "")}
                    required
                  />
                </Field>
              )}

              <Field>
                <FieldLabel htmlFor="password">Пароль</FieldLabel>

                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="Ваш пароль"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                  />

                  <button
                    type="button"
                    onMouseDown={(e) => e.preventDefault()}
                    onClick={() => setShowPassword((v) => !v)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground"
                  >
                    {showPassword ? (
                      <HugeiconsIcon icon={ViewOffIcon} size={18} />
                    ) : (
                      <HugeiconsIcon icon={ViewIcon} size={18} />
                    )}
                  </button>
                </div>
              </Field>

              <Field>
                <Button
                  type="submit"
                  className="w-full"
                  disabled={!email || !password}
                >
                  Войти
                </Button>

                <FieldDescription className="text-center">
                  Нет аккаунта?{" "}
                  <Link
                    href="/auth/signup"
                    className="underline underline-offset-4"
                  >
                    Зарегистрироваться
                  </Link>
                </FieldDescription>
              </Field>

              {successMessage && (
                <p className="text-center text-sm text-green-600">
                  {successMessage}
                </p>
              )}

              {errorMessage && (
                <p className="text-center text-sm text-red-500">
                  {errorMessage}
                </p>
              )}
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
