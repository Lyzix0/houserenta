import { cookies } from "next/headers"
import { redirect } from "next/navigation"

export default async function HomeLayout({
  children,
}: {
  children: React.ReactNode
}) {
  
  const bypassAuth = process.env.NEXT_PUBLIC_BYPASS_AUTH === "true" || process.env.NODE_ENV === "development";

  if (!bypassAuth) {
    const cookieStore = await cookies()
    // The backend's session middleware (fiber/v3 session) always names its cookie
    // "session_id" — this must match, or every visit here bounces straight back
    // to /auth/login even for an authenticated user.
    const sessionCookie = cookieStore.get("session_id")

    if (!sessionCookie) {
      redirect("/auth/login")
    }
  }

  return children
}
