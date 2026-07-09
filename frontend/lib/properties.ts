"use client";

import { apiFetch } from "./api";

async function handleResponse<T>(res: Response, errorMsg: string): Promise<T> {
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    let err: any = {};
    try { err = JSON.parse(text); } catch {}
    throw new Error(err.error || errorMsg);
  }
  return res.json() as Promise<T>;
}

export async function getProperties(): Promise<any[]> {
  const res = await apiFetch("/properties");
  return handleResponse<any[]>(res, "Не удалось получить список объектов.");
}

export async function createProperty(payload: {
  name: string;
  coordinates: string;
  country?: string;
  region: string;
  city: string;
  street: string;
  house: string;
  apartment: string;
  gvsTariff: number;
  hvsTariff: number;
  el1Tariff: number;
  el2Tariff?: number | null;
}): Promise<any> {
  const res = await apiFetch("/properties", {
    method: "POST",
    body: JSON.stringify(payload),
  });
  return handleResponse<any>(res, "Не удалось создать объект.");
}

export async function updateProperty(
  id: string,
  payload: {
    name: string;
    coordinates: string;
    country?: string;
    region: string;
    city: string;
    street: string;
    house: string;
    apartment: string;
    gvsTariff: number;
    hvsTariff: number;
    el1Tariff: number;
    el2Tariff?: number | null;
  }
): Promise<any> {
  const res = await apiFetch(`/properties/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
  return handleResponse<any>(res, "Не удалось обновить объект.");
}

export async function deleteProperty(id: string): Promise<any> {
  const res = await apiFetch(`/properties/${id}`, {
    method: "DELETE",
  });
  return handleResponse<any>(res, "Не удалось удалить объект.");
}

export async function createLease(
  propertyId: string,
  payload: {
    tenantUserId: string;
    price: number;
    monthsOfRent: number;
    paymentDay: number;
    readingDay: number;
  }
): Promise<any> {
  const res = await apiFetch(`/properties/${propertyId}/lease`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
  return handleResponse<any>(res, "Не удалось оформить аренду.");
}

export async function deleteLease(propertyId: string): Promise<any> {
  const res = await apiFetch(`/properties/${propertyId}/lease`, {
    method: "DELETE",
  });
  return handleResponse<any>(res, "Не удалось расторгнуть аренду.");
}

export async function submitReadings(
  propertyId: string,
  payload: {
    gvs: number;
    hvs: number;
    el1: number;
    el2?: number | null;
  }
): Promise<any> {
  const res = await apiFetch(`/properties/${propertyId}/readings`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
  return handleResponse<any>(res, "Не удалось передать показания счетчиков.");
}

export async function payBill(
  propertyId: string,
  payload: {
    amount: number;
    billId?: string;
  }
): Promise<any> {
  const res = await apiFetch(`/properties/${propertyId}/pay`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
  return handleResponse<any>(res, "Не удалось провести оплату.");
}

export async function addCustomItem(
  propertyId: string,
  payload: {
    description: string;
    amount: number;
  }
): Promise<any> {
  const res = await apiFetch(`/properties/${propertyId}/custom-item`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
  return handleResponse<any>(res, "Не удалось добавить начисление.");
}

export async function getVacantProperties(): Promise<any[]> {
  const res = await apiFetch("/properties/vacant");
  return handleResponse<any[]>(res, "Не удалось получить список свободных объектов.");
}

export async function applyForProperty(propertyId: string): Promise<any> {
  const res = await apiFetch(`/properties/${propertyId}/apply`, {
    method: "POST",
  });
  return handleResponse<any>(res, "Не удалось подать заявку на аренду.");
}

export async function getUnlinkedTenants(): Promise<any[]> {
  const res = await apiFetch("/tenants/unlinked");
  return handleResponse<any[]>(res, "Не удалось загрузить список арендаторов.");
}

export async function getChatHistory(propertyId: string): Promise<any[]> {
  const res = await apiFetch(`/chat/history/${propertyId}`);
  return handleResponse<any[]>(res, "Не удалось загрузить историю чата.");
}
