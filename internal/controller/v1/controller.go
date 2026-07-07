package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/potom_pridumaem/internal/repo/persistent"
	"github.com/potom_pridumaem/internal/usecase"
	"go.uber.org/zap"
)

type V1 struct {
	u            usecase.User
	p            usecase.Property
	propertyRepo *persistent.PropertyRepo
	l            *zap.Logger
	v            *validator.Validate
	sess         *session.Store
}
