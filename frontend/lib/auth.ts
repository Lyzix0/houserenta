"use client";

import { apiFetch } from "./api";

export type AuthUser = {
  id: string;
  name: string;
  email: string;
  role: "landlord" | "tenant" | "admin";
  allowedRoles: ("landlord" | "tenant" | "admin")[];
  document: string;
  phone: string;
  paymentCard?: string;
  tenantPropertyId?: string;
};

type RegisterPayload = {
  name: string;
  email: string;
  password: string;
  phone: string;
  initialRole: "landlord" | "tenant";
};

type LoginPayload = {
  email: string; 
  password: string;
};

type AuthResponse<T> = {
  ok: boolean;
  status: number;
  data: T;
};

async function parseResponse<T>(response: Response): Promise<AuthResponse<T>> {
  const text = await response.text();
  const data = text ? JSON.parse(text) : {};

  return {
    ok: response.ok,
    status: response.status,
    data,
  };
}

function getErrorMessage(error: unknown, fallback: string) {
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return fallback;
}

async function expectSuccess<T>(response: Response, fallback: string) {
  const parsed = await parseResponse<T & { error?: string }>(response);
  if (!parsed.ok) {
    throw new Error(parsed.data.error || fallback);
  }
  return parsed.data;
}

export function isUnauthorizedError(error: unknown) {
  return error instanceof Error && error.message === "Unauthorized";
}

export async function register(payload: RegisterPayload) {
  console.log("[register] payload:", payload);
  return await expectSuccess<{ user?: AuthUser; message?: string }>(
    await apiFetch("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
    "Регистрация не удалась.",
  );
}

export async function login(payload: LoginPayload) {
  console.log("[login] payload:", payload);
  return await expectSuccess<{ user?: AuthUser; message?: string }>(
    await apiFetch("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
    "Неверная почта или пароль.",
  );
}

export async function logout() {
  await apiFetch("/auth/logout", {
    method: "POST",
  });
}

export async function getMe(): Promise<AuthUser> {
  try {
    return await expectSuccess<AuthUser>(
      await apiFetch("/auth/me", {
        method: "GET",
      }),
      "Unauthorized",
    );
  } catch (error) {
    throw new Error(getErrorMessage(error, "Unauthorized"));
  }
}

export async function updateProfile(payload: {
  name: string;
  document: string;
  phone: string;
  paymentCard?: string;
  email: string;
  password?: string;
}) {
  return expectSuccess<{ ok: boolean }>(
    await apiFetch("/auth/profile", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
    "Не удалось обновить профиль.",
  );
}
