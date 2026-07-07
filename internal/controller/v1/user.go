package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/potom_pridumaem/internal/controller/middleware"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	"github.com/potom_pridumaem/internal/controller/v1/response"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/internal/usecase"
)

// register godoc
// @Summary      Register a user
// @Description  Creates a new user (landlord, tenant, or admin)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.Register  true  "Registration data"
// @Success      201    {object}  entity.User
// @Failure      400    {object}  response.Error  "invalid request body or role"
// @Failure      409    {object}  response.Error  "user already exists"
// @Failure      500    {object}  response.Error  "internal server error"
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
// @Summary      Log in a user
// @Description  Authenticates a user by email and password, and creates a session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.Login  true  "Login data"
// @Success      200    {object}  entity.User
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "invalid credentials"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /auth/login [post]
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
// @Summary      Log out a user
// @Description  Terminates the current user session
// @Tags         auth
// @Produce      json
// @Success      204  "session terminated successfully"
// @Failure      500  {object}  response.Error  "internal server error"
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
// @Summary      Current session
// @Description  Protected. Returns the full profile of the authenticated user for the current session, including the linked property when the caller is a tenant. Per business requirements this should also trigger an autobilling check; note: autobilling is not implemented yet in this API version, so only profile data is returned for now.
// @Tags         auth
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  response.Me
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /auth/me [get]
func (r *V1) me(ctx fiber.Ctx) error {
	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	profile, err := r.u.Me(ctx.Context(), userID)
	if err != nil {
		r.l.Error("restapi - v1 - me", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(response.Me{
		ID:               profile.User.ID,
		Name:             profile.User.Name,
		Email:            profile.User.Email,
		Role:             string(profile.User.Role),
		Document:         profile.User.Document,
		Phone:            profile.User.Phone,
		PaymentCard:      profile.User.PaymentCard,
		TenantPropertyID: profile.TenantPropertyID,
	})
}
