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
// @Description  Creates a new user (landlord or tenant) and immediately logs them in, setting the session cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.Register  true  "Registration data"
// @Success      200    {object}  response.AuthUser
// @Failure      400    {object}  response.Error  "invalid request body, invalid role, or email already registered"
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

	usr, err := r.u.Register(ctx.Context(), body.Name, body.Email, body.Password, body.InitialRole, body.Document, body.Phone, body.PaymentCard)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrEmailAlreadyTaken):
			return errorResponse(ctx, http.StatusBadRequest, "Пользователь с такой почтой уже зарегистрирован")
		case errors.Is(err, entity.ErrInvalidRole):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())
		default:
			r.l.Error("restapi - v1 - register", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
	}

	if err := r.startSession(ctx, usr); err != nil {
		r.l.Error("restapi - v1 - register", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(response.AuthUser{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email,
		Role:  string(usr.Role),
	})
}

// login godoc
// @Summary      Log in a user
// @Description  Authenticates a user by email or phone and password, and creates a session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.Login  true  "Login data"
// @Success      200    {object}  response.AuthUser
// @Failure      400    {object}  response.Error  "invalid request body or credentials"
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
			return errorResponse(ctx, http.StatusBadRequest, "Неверная почта или пароль")
		}
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	if err := r.startSession(ctx, usr); err != nil {
		r.l.Error("restapi - v1 - login", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(response.AuthUser{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email,
		Role:  string(usr.Role),
	})
}

// logout godoc
// @Summary      Log out a user
// @Description  Terminates the current user session
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]bool
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

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// startSession regenerates the session and stores the caller's identity in it,
// shared by register (auto-login) and login.
func (r *V1) startSession(ctx fiber.Ctx, usr entity.User) error {
	sess, err := r.sess.Get(ctx)
	if err != nil {
		return err
	}
	defer sess.Release()

	if err := sess.Regenerate(); err != nil {
		return err
	}

	sess.Set(middleware.UserIDLocalsKey, usr.ID)
	sess.Set(middleware.UserRoleLocalsKey, string(usr.Role))

	return sess.Save()
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

// profile godoc
// @Summary      Update account settings
// @Description  Protected. Partially updates the authenticated user's personal data, including changing the password. Only the provided fields are changed.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        input  body      request.Profile  true  "Profile fields to update"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body or email already registered"
// @Failure      401    {object}  response.Error  "not authenticated"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /auth/profile [post]
func (r *V1) profile(ctx fiber.Ctx) error {
	var body request.Profile
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("json update profile:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("validate profile:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	if err := r.u.UpdateProfile(ctx.Context(), userID, body); err != nil {
		switch {
		case errors.Is(err, repo.ErrEmailAlreadyTaken):
			return errorResponse(ctx, http.StatusBadRequest, "Пользователь с такой почтой уже зарегистрирован")
		default:
			r.l.Error("restapi - v1 - profile", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}
