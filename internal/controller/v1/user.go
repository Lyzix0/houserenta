package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/potom_pridumaem/internal/controller/middleware"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/internal/usecase"
)

func (r *V1) register(ctx fiber.Ctx) error {
	var body request.Register
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("restapi - v1 - register", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("restapi - v1 - register", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	usr, err := r.u.Register(ctx.Context(), body.Name, body.Email, body.Password, body.Role, body.Document, body.Phone)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrEmailAlreadyTaken):
			return errorResponse(ctx, http.StatusConflict, "user already exists")
		case errors.Is(err, entity.ErrInvalidRole):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())
		default:
			r.l.Error("restapi - v1 - register", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
	}

	return ctx.Status(http.StatusCreated).JSON(usr)
}

func (r *V1) login(ctx fiber.Ctx) error {
	var body request.Login
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	usr, err := r.u.Login(ctx.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return errorResponse(ctx, http.StatusUnauthorized, "invalid credentials")
		}
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	sess, err := r.sess.Get(ctx)
	if err != nil {
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	defer sess.Release()

	if err := sess.Regenerate(); err != nil {
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	sess.Set(middleware.UserIDLocalsKey, usr.ID)
	sess.Set(middleware.UserRoleLocalsKey, string(usr.Role))

	if err := sess.Save(); err != nil {
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(usr)
}

func (r *V1) logout(ctx fiber.Ctx) error {
	sess, err := r.sess.Get(ctx)
	if err != nil {
		r.l.Error("restapi - v1 - logout", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	defer sess.Release()

	if err := sess.Destroy(); err != nil {
		r.l.Error("restapi - v1 - logout", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (r *V1) me(ctx fiber.Ctx) error {
	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)
	role, _ := ctx.Locals(middleware.UserRoleLocalsKey).(string)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"id":   userID,
		"role": role,
	})
}
