package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

const (
	UserIDLocalsKey   = "user_id"
	UserRoleLocalsKey = "user_role"
)

func AuthRequired(store *session.Store) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		defer sess.Release()

		userID, ok := sess.Get(UserIDLocalsKey).(string)
		if !ok || userID == "" {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		ctx.Locals(UserIDLocalsKey, userID)
		if role, ok := sess.Get(UserRoleLocalsKey).(string); ok {
			ctx.Locals(UserRoleLocalsKey, role)
		}

		return ctx.Next()
	}
}
