"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { HugeiconsIcon } from "@hugeicons/react";
import { ViewIcon, ViewOffIcon } from "@hugeicons/core-free-icons";
import { register } from "@/lib/auth";

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

export function SignupForm({
  ...props
}: React.ComponentProps<typeof Card>) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<"landlord" | "tenant">(
    (searchParams.get("role") === "tenant" ? "tenant" : "landlord")
  );

  const [showPassword, setShowPassword] = useState(false);
  const [showHint, setShowHint] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const passwordRules = [
    {
      label: "Минимум 8 символов",
      check: (v: string) => v.length >= 8,
    },
    {
      label: "Хотя бы одна латинская буква",
      check: (v: string) => /[a-zA-Z]/.test(v),
    },
    {
      label: "Без пробелов",
      check: (v: string) => !/\s/.test(v),
    },
  ];

  const handleRegister = async (
    e: React.FormEvent<HTMLFormElement>
  ) => {
    e.preventDefault();
    setErrorMessage("");

    const isPasswordValid = passwordRules.every(rule => rule.check(password));
    if (!isPasswordValid) {
      setErrorMessage("Пароль не соответствует требованиям безопасности.");
      return;
    }

    try {
      await register({
        name,
        email,
        password,
        phone,
        initialRole: role,
      });

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
    <Card {...props}>
      <CardHeader>
        <CardTitle className="text-xl">Регистрация</CardTitle>
        <CardDescription>
          Создайте единый аккаунт в нашей экосистеме
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleRegister}>
          <FieldGroup className="gap-4">
            <Field>
              <FieldLabel>Выберите вашу основную роль *</FieldLabel>
              <p className="text-[10px] text-muted-foreground">Вы сможете переключать роли в любой момент в профиле.</p>
              <div className="grid grid-cols-2 mt-1">
                <Button
                  type="button"
                  variant={role === "landlord" ? "default" : "outline"}
                  onClick={() => setRole("landlord")}
                  className={`text-xs font-bold transition duration-150 rounded-r-none rounded-l-full ${
                    role === "landlord"
                      ? "border-0"
                      : "bg-slate-50 text-slate-700 hover:bg-slate-100"
                  }`}
                >
                  Арендодатель
                </Button>
                <Button
                  type="button"
                  variant={role === "tenant" ? "default" : "outline"}
                  onClick={() => setRole("tenant")}
                  className={`text-xs font-bold transition duration-150 rounded-l-none rounded-r-full ${
                    role === "tenant"
                      ? "border-0"
                      : "bg-slate-50  text-slate-700 hover:bg-slate-100"
                  }`}
                >
                  Арендатор
                </Button>
              </div>
            </Field>

            <Field>
              <FieldLabel htmlFor="name">Полное Имя</FieldLabel>
              <Input
                id="name"
                placeholder="Иванов Иван"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </Field>

            <Field>
              <FieldLabel htmlFor="phone">Телефон</FieldLabel>
              <PhoneInput
                id="phone"
                placeholder="+7 900 000 00 00"
                value={phone}
                onChange={(val) => setPhone(val || "")}
                required
              />
            </Field>

            <Field>
              <FieldLabel htmlFor="email">Почта</FieldLabel>
              <Input
                id="email"
                type="email"
                placeholder="landlord@mail.ru"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </Field>

            <Field>
              <FieldLabel htmlFor="password">Пароль</FieldLabel>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  placeholder="Придумайте надежный пароль"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  onFocus={() => setShowHint(true)}
                  onBlur={() => setShowHint(false)}
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

              <div
                className={`overflow-hidden transition-all duration-300 ${
                  showHint
                    ? "max-h-32 opacity-100"
                    : "max-h-0 opacity-0"
                }`}
              >
                <div className="mt-2 rounded-2xl border bg-muted p-3 text-xs">
                  <ul className="space-y-1">
                    {passwordRules.map((rule) => {
                      const passed = rule.check(password);

                      return (
                        <li
                          key={rule.label}
                          className={`flex items-center gap-2 ${
                            passed
                              ? "text-green-600"
                              : "text-red-500"
                          }`}
                        >
                          <span>{passed ? "✓" : "•"}</span>
                          {rule.label}
                        </li>
                      );
                    })}
                  </ul>
                </div>
              </div>
            </Field>

            {errorMessage && (
              <p className="text-center text-sm text-red-500">
                {errorMessage}
              </p>
            )}

            <Field className="pt-2">
              <Button
                type="submit"
                className="w-full"
                disabled={!name || !phone || !email || !password}
              >
                Зарегистрироваться
              </Button>

              <FieldDescription className="text-center">
                Уже есть аккаунт?{" "}
                <Link
                  href="/auth/login"
                  className="underline underline-offset-4"
                >
                  Войти
                </Link>
              </FieldDescription>
            </Field>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
}
