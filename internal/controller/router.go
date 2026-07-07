package restapi

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/potom_pridumaem/config"
	"github.com/potom_pridumaem/internal/controller/middleware"
	v1 "github.com/potom_pridumaem/internal/controller/v1"
	"github.com/potom_pridumaem/internal/repo/persistent"
	"github.com/potom_pridumaem/internal/usecase"
)

func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	u usecase.User,
	p usecase.Property,
	propertyRepo *persistent.PropertyRepo,
	l *zap.Logger,
) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// проверка на живучесть сайта)
	app.Get("/healthz", func(ctx fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})

	apiV1Group := app.Group("/v1")
	v1.NewRoutes(apiV1Group, u, p, l, propertyRepo)
}
