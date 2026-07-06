package middleware

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v3"
	fiberRecover "github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
)

func buildPanicMessage(ctx fiber.Ctx, err any) string {
	var result strings.Builder
	result.WriteString(ctx.IP())
	result.WriteString(" - ")
	result.WriteString(ctx.Method())
	result.WriteString(" ")
	result.WriteString(ctx.OriginalURL())
	result.WriteString(" PANIC DETECTED: ")
	fmt.Fprintf(&result, "%v\n%s\n", err, debug.Stack())
	return result.String()
}

func logPanic(l *zap.Logger) func(c fiber.Ctx, err any) {
	return func(ctx fiber.Ctx, err any) {
		l.Error(buildPanicMessage(ctx, err))
	}
}

func Recovery(l *zap.Logger) func(fiber.Ctx) error {
	return fiberRecover.New(fiberRecover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: logPanic(l),
	})
}
