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
	p usecase.Property, l *zap.Logger,
) {
	sess := session.NewStore(session.Config{
		IdleTimeout:    24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	r := &V1{
		p:    p,
		u:    u,
		l:    l,
		v:    validator.New(validator.WithRequiredStructEnabled()),
		sess: sess,
	}

	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.Post("/login", r.login)
		authGroup.Post("/register", r.register)
		authGroup.Post("/logout", r.logout)
		authGroup.Get("/me", middleware.AuthRequired(sess), r.me)
	}

	propertyGroup := apiV1Group.Group("/properties")
	{
		propertyGroup.Post("/property", r.createProperty)
		propertyGroup.Get("/", middleware.AuthRequired(sess), r.getProperties)
		propertyGroup.Get("/:id", middleware.AuthRequired(sess), r.getProperty)
		propertyGroup.Put("/:id", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.updateProperty)
		propertyGroup.Delete("/:id", middleware.AuthRequired(sess), middleware.RoleRequired(string(entity.RoleLandlord)), r.deleteProperty)
	}
}
