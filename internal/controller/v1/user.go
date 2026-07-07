package v1

import (
	"errors"
	"fmt"
	"math/rand"
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

	if body.PaymentCard != nil && *body.PaymentCard == "" {
		body.PaymentCard = nil
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

	sess, err := r.sess.Get(ctx)
	if err != nil {
		r.l.Error("restapi - v1 - register - session", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	defer sess.Release()

	if err := sess.Regenerate(); err != nil {
		r.l.Error("restapi - v1 - register - regenerate", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	sess.Set(middleware.UserIDLocalsKey, usr.ID)
	sess.Set(middleware.UserRoleLocalsKey, string(usr.Role))

	if err := sess.Save(); err != nil {
		r.l.Error("restapi - v1 - register - save", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
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
// @Summary      Current user
// @Description  Returns the user ID and role from the current session
// @Tags         auth
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  response.Error  "not authenticated"
// @Router       /auth/me [get]
func (r *V1) me(ctx fiber.Ctx) error {
	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	usr, err := r.u.GetByID(ctx.Context(), userID)
	if err != nil {
		r.l.Error("restapi - v1 - me", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"id":       usr.ID,
		"name":     usr.Name,
		"email":    usr.Email,
		"role":     string(usr.Role),
		"document": usr.Document,
		"phone":    usr.Phone,
	})
}

// profile godoc
// @Summary      Update user profile
// @Description  Updates name, document, phone, email and payment card
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        input  body      request.UpdateProfile  false  "Profile fields to update"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /auth/profile [post]
func (r *V1) profile(ctx fiber.Ctx) error {
	userID := middlewareGetUserID(ctx)
	var body request.UpdateProfile
	if err := ctx.Bind().Body(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := r.u.UpdateProfile(ctx.Context(), userID, body.Name, body.Document, body.Phone, body.Email, body.PaymentCard); err != nil {
		r.l.Error("restapi - v1 - profile", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// changePassword godoc
// @Summary      Change user password
// @Description  Changes password after verifying the old one
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        input  body      request.ChangePassword  true  "Old and new passwords"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "invalid old password"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /auth/change-password [post]
func (r *V1) changePassword(ctx fiber.Ctx) error {
	userID := middlewareGetUserID(ctx)
	var body request.ChangePassword
	if err := ctx.Bind().Body(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := r.u.ChangePassword(ctx.Context(), userID, body.OldPassword, body.NewPassword); err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return errorResponse(ctx, http.StatusUnauthorized, "invalid old password")
		}
		r.l.Error("restapi - v1 - changePassword", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// switchRole godoc
// @Summary      Switch active user role
// @Description  Changes current user role to another valid role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        input  body      request.SwitchRole  true  "Target role"
// @Success      200    {object}  entity.User
// @Failure      400    {object}  response.Error  "invalid request body or role"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /auth/switch-role [post]
func (r *V1) switchRole(ctx fiber.Ctx) error {
	userID := middlewareGetUserID(ctx)
	var body request.SwitchRole
	if err := ctx.Bind().Body(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	usr, err := r.u.SwitchRole(ctx.Context(), userID, body.TargetRole)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidRole) {
			return errorResponse(ctx, http.StatusBadRequest, "invalid role")
		}
		r.l.Error("restapi - v1 - switchRole", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"id":    usr.ID,
		"name":  usr.Name,
		"email": usr.Email,
		"role":  string(usr.Role),
	})
}

// forgotPassword godoc
// @Summary      Request password reset code
// @Description  Sends a mock password reset code to the given email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.ForgotPassword  true  "User email"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  response.Error  "invalid request body"
// @Router       /auth/forgot-password [post]
func (r *V1) forgotPassword(ctx fiber.Ctx) error {
	var body request.ForgotPassword
	if err := ctx.Bind().Body(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	user, err := r.u.GetByEmail(ctx.Context(), body.Email)
	if err != nil {
		r.l.Debug("forgotPassword - user not found", zap.String("email", body.Email))
		return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "If the email exists, a reset code was sent"})
	}
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	r.l.Info(fmt.Sprintf("[MAIL MOCK] Password reset code for %s: %s", user.Email, code))
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "If the email exists, a reset code was sent"})
}

// resetPassword godoc
// @Summary      Reset password
// @Description  Resets user password with a valid code (mock)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.ResetPassword  true  "Email, code and new password"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  response.Error  "invalid request body"
// @Router       /auth/reset-password [post]
func (r *V1) resetPassword(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Password changed (mock)"})
}

// requestEmailVerify godoc
// @Summary      Request email verification code
// @Description  Sends a mock verification code to the given email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.ForgotPassword  true  "User email"
// @Success      200    {object}  map[string]string
// @Router       /auth/request-email-verify [post]
func (r *V1) requestEmailVerify(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Verification code sent (mock)"})
}

// verifyEmail godoc
// @Summary      Verify email with code
// @Description  Verifies user email with a valid code (mock)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      request.VerifyEmail  true  "Email and code"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  response.Error  "invalid code"
// @Router       /auth/verify-email [post]
func (r *V1) verifyEmail(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Email verified (mock)"})
}

func middlewareGetUserID(ctx fiber.Ctx) string {
	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)
	return userID
}
