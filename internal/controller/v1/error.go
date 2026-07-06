package v1

import (
	"github.com/gofiber/fiber/v3"
	"github.com/potom_pridumaem/internal/controller/v1/response"
)

func errorResponse(ctx fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Error{Error: msg})
}
