"use client"

import { useQuery } from "@tanstack/react-query"

import { getMe } from "@/lib/auth"

export function useUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: getMe,
    staleTime: 0,
    retry: false,
  })
}
