package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func Logger(l *zap.Logger) func(c fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		err := ctx.Next()

		status := ctx.Response().StatusCode()

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			status = fiberErr.Code
		}

		fields := []zap.Field{
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.Int("status", status),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			l.Error("request failed", fields...)
			return err
		}

		l.Info("request completed", fields...)
		return err
	}
}
