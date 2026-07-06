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

// register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя (арендодателя, арендатора или администратора)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.Register  true  "Данные для регистрации"
// @Success      201    {object}  entity.User
// @Failure      400    {object}  response.ErrorResponse  "невалидное тело запроса или роль"
// @Failure      409    {object}  response.ErrorResponse  "пользователь уже существует"
// @Failure      500    {object}  response.ErrorResponse  "внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (r *V1) register(ctx fiber.Ctx) error {
	var body request.Register
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("json register user:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("register user:", zap.Error(err))
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

// login godoc
// @Summary      Вход пользователя
// @Description  Аутентифицирует пользователя по email и паролю, создаёт сессию
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.Login  true  "Данные для входа"
// @Success      200    {object}  entity.User
// @Failure      400    {object}  response.ErrorResponse  "невалидное тело запроса"
// @Failure      401    {object}  response.ErrorResponse  "неверные учётные данные"
// @Failure      500    {object}  response.ErrorResponse  "внутренняя ошибка сервера"
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

// logout godoc
// @Summary      Выход пользователя
// @Description  Завершает текущую сессию пользователя
// @Tags         auth
// @Produce      json
// @Success      204  "сессия успешно завершена"
// @Failure      500  {object}  response.ErrorResponse  "внутренняя ошибка сервера"
// @Router       /auth/logout [post]
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

// me godoc
// @Summary      Текущий пользователь
// @Description  Возвращает ID и роль пользователя из текущей сессии
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /auth/me [get]
func (r *V1) me(ctx fiber.Ctx) error {
	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)
	role, _ := ctx.Locals(middleware.UserRoleLocalsKey).(string)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"id":   userID,
		"role": role,
	})
}
