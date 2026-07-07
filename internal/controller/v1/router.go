package v1

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/potom_pridumaem/internal/controller/middleware"
	"github.com/potom_pridumaem/internal/repo/persistent"
	"github.com/potom_pridumaem/internal/usecase"
	"go.uber.org/zap"
)

func NewRoutes(
	apiV1Group fiber.Router, u usecase.User,
	p usecase.Property, l *zap.Logger,
	propertyRepo *persistent.PropertyRepo,
) {
	sess := session.NewStore(session.Config{
		IdleTimeout:    24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	r := &V1{
		p:            p,
		u:            u,
		propertyRepo: propertyRepo,
		l:            l,
		v:            validator.New(validator.WithRequiredStructEnabled()),
		sess:         sess,
	}

	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.Post("/login", r.login)
		authGroup.Post("/register", r.register)
		authGroup.Post("/logout", r.logout)
		authGroup.Get("/me", middleware.AuthRequired(sess), r.me)
		authGroup.Post("/profile", middleware.AuthRequired(sess), r.profile)
		authGroup.Post("/change-password", middleware.AuthRequired(sess), r.changePassword)
		authGroup.Post("/switch-role", middleware.AuthRequired(sess), r.switchRole)
		authGroup.Post("/forgot-password", r.forgotPassword)
		authGroup.Post("/reset-password", r.resetPassword)
		authGroup.Post("/request-email-verify", r.requestEmailVerify)
		authGroup.Post("/verify-email", r.verifyEmail)
	}

	propertyGroup := apiV1Group.Group("/properties", middleware.AuthRequired(sess))
	{
		propertyGroup.Post("/property", r.createProperty)
		propertyGroup.Get("/", r.getProperties)
		propertyGroup.Get("", r.getProperties)
		propertyGroup.Get("/:id", r.getProperty)
		propertyGroup.Put("/:id", r.updateProperty)
		propertyGroup.Delete("/:id", r.deleteProperty)
		propertyGroup.Get("/vacant", r.getVacantProperties)
	}

	tenantsGroup := apiV1Group.Group("/tenants", middleware.AuthRequired(sess))
	{
		tenantsGroup.Get("/unlinked", r.getUnlinkedTenants)
	}
}
