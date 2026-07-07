package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/potom_pridumaem/internal/controller/middleware"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	"github.com/potom_pridumaem/internal/controller/v1/response"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"go.uber.org/zap"
)

// createProperty godoc
// @Summary      Create a property
// @Description  Protected, Landlord only. Creates a new rental property owned by the authenticated landlord, registering its address and utility tariffs (hot/cold water, electricity) that will later be used for billing calculations. The new property starts with a zero balance.
// @Tags         properties
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        input  body      request.Property  true  "Property data"
// @Success      200    {object}  response.PropertySummary
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "not authenticated"
// @Failure      404    {object}  response.Error  "landlord not found"
// @Failure      409    {object}  response.Error  "property already exists"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties [post]
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

	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	prop, err := r.p.CreateProperty(ctx.Context(), landlordID, body)
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

	return ctx.Status(http.StatusOK).JSON(response.PropertySummary{
		ID:   prop.ID,
		Name: prop.Name,
	})
}

// getProperties godoc
// @Summary      List properties
// @Description  Protected. Returns the caller's properties with an automatic nested assembly of related entities (meter readings, bills with their line items, upcoming custom charges, lease/tenant info, landlord contact). Landlords see every property they own; tenants see the single property they lease, if any. The "applications" field is always empty: there is no applications/tickets subsystem implemented yet.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Success      200  {array}   entity.PropertyDetail
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties [get]
func (r *V1) getProperties(ctx fiber.Ctx) error {
	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)
	role, _ := ctx.Locals(middleware.UserRoleLocalsKey).(string)

	properties, err := r.p.GetProperties(ctx.Context(), userID, entity.Role(role))
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

// createLease godoc
// @Summary      Move a tenant in
// @Description  Protected, Landlord only. Creates a lease linking a registered tenant to the property. If the property already has an active lease, it is superseded by this one. Per business requirements this should also automatically cancel other third-party applications for the property; note: there is no applications/tickets subsystem implemented yet, so this step is a no-op.
// @Tags         properties
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id     path      string         true  "Property ID"
// @Param        input  body      request.Lease  true  "Lease data"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "not authenticated"
// @Failure      404    {object}  response.Error  "property not found, or tenant not found"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties/{id}/lease [post]
func (r *V1) createLease(ctx fiber.Ctx) error {
	var body request.Lease
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("create lease json:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("validate lease:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	if err := r.p.CreateLease(ctx.Context(), ctx.Params("id"), landlordID, body); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		case errors.Is(err, repo.ErrTenantNotFound):
			return errorResponse(ctx, http.StatusNotFound, "Жилец не найден в системе")
		default:
			r.l.Error("create lease:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to create lease")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// deleteLease godoc
// @Summary      Move a tenant out
// @Description  Protected, Landlord only. Terminates the property's lease and unlinks the tenant. Idempotent: succeeds even if the property currently has no lease.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Param        id   path      string  true  "Property ID"
// @Success      200  {object}  map[string]bool
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      404  {object}  response.Error  "property not found"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/{id}/lease [delete]
func (r *V1) deleteLease(ctx fiber.Ctx) error {
	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	if err := r.p.DeleteLease(ctx.Context(), ctx.Params("id"), landlordID); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		default:
			r.l.Error("delete lease:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to delete lease")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// deleteProperty godoc
// @Summary      Delete a property
// @Description  Protected, Landlord only. Permanently deletes a property. Any lease bound to it is automatically removed as well (database cascade), unlinking the tenant. Only the owning landlord may delete their property.
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
