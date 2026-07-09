package v1

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/potom_pridumaem/internal/controller/middleware"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/usecase"
	"go.uber.org/zap"
)

func NewRoutes(
	apiV1Group fiber.Router, u usecase.User,
	p usecase.Property, billing usecase.Billing, l *zap.Logger,
) {
	sess := session.NewStore(session.Config{
		IdleTimeout:    24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		// The public entry point is now Caddy terminating HTTPS on the real domain;
		// browsers refuse to send a Secure cookie back over a plain-HTTP origin, so
		// local testing must go through the domain (or drop this when testing bare HTTP).
		CookieSecure: true,
	})

	r := &V1{
		p:       p,
		u:       u,
		billing: billing,
		l:       l,
		v:       validator.New(validator.WithRequiredStructEnabled()),
		sess:    sess,
	}

	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.Post("/login", r.login)
		authGroup.Post("/register", r.register)
		authGroup.Post("/logout", r.logout)
		authGroup.Get("/me", middleware.AuthRequired(sess), r.me)
		authGroup.Post("/profile", middleware.AuthRequired(sess), r.profile)
	}

	propertyGroup := apiV1Group.Group("/properties")
	{
		propertyGroup.Post("/", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.createProperty)
		propertyGroup.Get("/", middleware.AuthRequired(sess), r.getProperties)
		propertyGroup.Get("/vacant", middleware.AuthRequired(sess), r.getVacantProperties)
		propertyGroup.Get("/:id", middleware.AuthRequired(sess), r.getProperty)
		propertyGroup.Put("/:id", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.updateProperty)
		propertyGroup.Delete("/:id", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.deleteProperty)
		propertyGroup.Post("/:id/lease", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.createLease)
		propertyGroup.Delete("/:id/lease", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.deleteLease)
		propertyGroup.Post("/:id/readings", middleware.AuthRequired(sess), r.createReading)
		propertyGroup.Post("/:id/pay", middleware.AuthRequired(sess), r.pay)
		propertyGroup.Post("/:id/custom-item", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.createCustomItem)
		propertyGroup.Post("/:id/apply", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleTenant)), r.applyToProperty)
	}

	tenantGroup := apiV1Group.Group("/tenants")
	{
		tenantGroup.Get("/unlinked", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.getUnlinkedTenants)
	}
}
