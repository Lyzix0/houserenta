"use client"

import { ReactNode, useEffect, useState } from "react"
import { usePathname, useRouter } from "next/navigation"

import { getMe, isUnauthorizedError, type AuthUser } from "@/lib/auth"

type RequireAuthProps = {
  children: ReactNode
}

export function RequireAuth({ children }: RequireAuthProps) {
  const router = useRouter()
  const pathname = usePathname()
  const [user, setUser] = useState<AuthUser | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    let active = true

    async function loadUser() {
      try {
        const currentUser = await getMe()

        if (!active) {
          return
        }

        setUser(currentUser)
      } catch (authError) {
        if (!active) {
          return
        }

        if (isUnauthorizedError(authError)) {
          const next = pathname ? `?next=${encodeURIComponent(pathname)}` : ""
          router.replace(`/auth/login${next}`)
          return
        }

        setError("Не удалось проверить авторизацию.")
      } finally {
        if (active) {
          setLoading(false)
        }
      }
    }

    loadUser()

    return () => {
      active = false
    }
  }, [pathname, router])

  if (loading) {
    return <div className="p-4">Проверяем авторизацию...</div>
  }

  if (error) {
    return <div className="p-4 text-red-500">{error}</div>
  }

  if (!user) {
    return null
  }

  return <>{children}</>
}
