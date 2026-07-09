"use client"

import { useQuery } from "@tanstack/react-query"
import { getProperties, getVacantProperties, getUnlinkedTenants } from "@/lib/properties"
import { useUser } from "./use-user"

interface QueryOptions {
  enabled?: boolean
}

export function useProperties(options?: QueryOptions) {
  const { data: user } = useUser()
  return useQuery({
    queryKey: ["properties", user?.id],
    queryFn: getProperties,
    staleTime: 30 * 1000,
    refetchInterval: 30 * 1000,
    enabled: options?.enabled,
  })
}

export function useVacantProperties(options?: QueryOptions) {
  const { data: user } = useUser()
  return useQuery({
    queryKey: ["vacantProperties", user?.id],
    queryFn: getVacantProperties,
    staleTime: 30 * 1000,
    refetchInterval: 30 * 1000,
    enabled: options?.enabled,
  })
}

export function useUnlinkedTenants(options?: QueryOptions) {
  const { data: user } = useUser()
  return useQuery({
    queryKey: ["unlinkedTenants", user?.id],
    queryFn: getUnlinkedTenants,
    staleTime: 30 * 1000,
    refetchInterval: 30 * 1000,
    enabled: options?.enabled,
  })
}
