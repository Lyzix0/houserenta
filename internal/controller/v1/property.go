package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/potom_pridumaem/internal/controller/middleware"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	"github.com/potom_pridumaem/internal/repo"
	"go.uber.org/zap"
)

// createProperty godoc
// @Summary      Create a property
// @Description  Creates a new rental property for a landlord, registering its address and utility tariffs (hot/cold water, electricity) that will later be used for billing calculations. The new property starts with a zero balance.
// @Tags         properties
// @Accept       json
// @Produce      json
// @Param        input  body      request.Property  true  "Property data"
// @Success      201    {object}  entity.Property
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      404    {object}  response.Error  "landlord not found"
// @Failure      409    {object}  response.Error  "property already exists"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties/property [post]
func (r *V1) createProperty(ctx fiber.Ctx) error {
	var body request.Property
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("create property json:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("validate property:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	prop, err := r.p.CreateProperty(ctx.Context(), body)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrLandlordNotFound):
			return errorResponse(ctx, http.StatusNotFound, "landlord not found")
		case errors.Is(err, repo.ErrPropertyAlreadyExists):
			return errorResponse(ctx, http.StatusConflict, "property already exists")
		default:
			r.l.Error("create property:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to create property")
		}
	}

	return ctx.Status(http.StatusCreated).JSON(prop)
}

// getProperties godoc
// @Summary      List properties
// @Description  Returns all properties owned by the currently authenticated user, so a landlord can see their full property portfolio. Only properties the caller owns are returned; other landlords' properties are never exposed.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Success      200  {array}   entity.Property
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties [get]
func (r *V1) getProperties(ctx fiber.Ctx) error {
	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	properties, err := r.p.GetProperties(ctx.Context(), landlordID)
	if err != nil {
		r.l.Error("get properties:", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get properties")
	}

	return ctx.Status(http.StatusOK).JSON(properties)
}

// getProperty godoc
// @Summary      Get a property
// @Description  Returns a single property by ID. The property must belong to the authenticated landlord; a property owned by another landlord is reported as not found, so ownership can never be probed via this endpoint.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Param        id   path      string  true  "Property ID"
// @Success      200  {object}  entity.Property
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      404  {object}  response.Error  "property not found"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/{id} [get]
func (r *V1) getProperty(ctx fiber.Ctx) error {
	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	prop, err := r.p.GetProperty(ctx.Context(), ctx.Params("id"), landlordID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		default:
			r.l.Error("get property:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to get property")
		}
	}

	return ctx.Status(http.StatusOK).JSON(prop)
}

// updateProperty godoc
// @Summary      Update a property
// @Description  Protected, Landlord only. Updates a property's characteristics and utility tariffs (address, area details, hot/cold water and electricity tariffs). Accepts the same payload shape as property creation. Only the owning landlord may edit their property.
// @Tags         properties
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id     path      string            true  "Property ID"
// @Param        input  body      request.Property  true  "Property data"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "not authenticated"
// @Failure      404    {object}  response.Error  "property not found"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties/{id} [put]
func (r *V1) updateProperty(ctx fiber.Ctx) error {
	var body request.Property
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("update property json:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("validate property:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	if err := r.p.UpdateProperty(ctx.Context(), ctx.Params("id"), landlordID, body); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		default:
			r.l.Error("update property:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to update property")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// deleteProperty godoc
// @Summary      Delete a property
// @Description  Protected, Landlord only. Permanently deletes a property. Per business requirements this should also automatically terminate any linked rental contracts and clear the rental reference for bound tenants; note that contract/tenancy linkage does not exist yet in this API version, so only the property record itself is removed for now. Only the owning landlord may delete their property.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Param        id   path      string  true  "Property ID"
// @Success      200  {object}  map[string]bool
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      404  {object}  response.Error  "property not found"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/{id} [delete]
func (r *V1) deleteProperty(ctx fiber.Ctx) error {
	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	if err := r.p.DeleteProperty(ctx.Context(), ctx.Params("id"), landlordID); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		default:
			r.l.Error("delete property:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to delete property")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}
