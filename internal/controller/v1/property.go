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
// @Description  Protected. Returns the caller's properties with an automatic nested assembly of related entities (meter readings, bills with their line items, upcoming custom charges, lease/tenant info, landlord contact). Landlords see every property they own; tenants see the single property they lease, if any. The "applications" field is always empty: there is no applications/tickets subsystem implemented yet. Also triggers the lazy auto-billing check: for each active lease, if 30+ days have passed since its last rent bill (or none exists yet), a new bill is generated from the lease price, any queued custom charges, and unbilled utility consumption; a billing failure is logged and does not affect this response.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Success      200  {array}   entity.PropertyDetail
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties [get]
func (r *V1) getProperties(ctx fiber.Ctx) error {
	r.runBillingCheck(ctx)

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

// createReading godoc
// @Summary      Submit a meter reading
// @Description  Protected. Records a new meter reading for the property, submitted by the tenant currently leasing it or by the owning landlord.
// @Tags         properties
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id     path      string           true  "Property ID"
// @Param        input  body      request.Reading  true  "Meter reading"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "not authenticated"
// @Failure      404    {object}  response.Error  "property not found"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties/{id}/readings [post]
func (r *V1) createReading(ctx fiber.Ctx) error {
	var body request.Reading
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("create reading json:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("validate reading:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)
	role, _ := ctx.Locals(middleware.UserRoleLocalsKey).(string)

	if err := r.p.CreateReading(ctx.Context(), ctx.Params("id"), userID, entity.Role(role), body); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		default:
			r.l.Error("create reading:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to submit reading")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// pay godoc
// @Summary      Make a payment
// @Description  Protected. Either settles a specific bill (status becomes "paid") when billId is provided, or tops up the property's general balance otherwise.
// @Tags         properties
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id     path      string           true  "Property ID"
// @Param        input  body      request.Payment  true  "Payment data"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "not authenticated"
// @Failure      404    {object}  response.Error  "property not found, or bill not found"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties/{id}/pay [post]
func (r *V1) pay(ctx fiber.Ctx) error {
	var body request.Payment
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("pay json:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("validate payment:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	userID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)
	role, _ := ctx.Locals(middleware.UserRoleLocalsKey).(string)

	if err := r.p.Pay(ctx.Context(), ctx.Params("id"), userID, entity.Role(role), body); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		case errors.Is(err, repo.ErrBillNotFound):
			return errorResponse(ctx, http.StatusNotFound, "bill not found")
		default:
			r.l.Error("pay:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to process payment")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// createCustomItem godoc
// @Summary      Add a one-off charge
// @Description  Protected, Landlord only. Adds a one-off charge (e.g. cleaning, lock replacement) that will be folded into the property's next automatic bill.
// @Tags         properties
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id     path      string              true  "Property ID"
// @Param        input  body      request.CustomItem  true  "Charge data"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      401    {object}  response.Error  "not authenticated"
// @Failure      404    {object}  response.Error  "property not found"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties/{id}/custom-item [post]
func (r *V1) createCustomItem(ctx fiber.Ctx) error {
	var body request.CustomItem
	if err := ctx.Bind().Body(&body); err != nil {
		r.l.Error("create custom item json:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error("validate custom item:", zap.Error(err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	landlordID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	if err := r.p.CreateCustomItem(ctx.Context(), ctx.Params("id"), landlordID, body); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		default:
			r.l.Error("create custom item:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to create custom item")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// runBillingCheck performs the lazy auto-billing pass. It is a passive side effect
// of calling this route, not the route's own concern, so a failure here is logged
// and swallowed rather than failing the request.
func (r *V1) runBillingCheck(ctx fiber.Ctx) {
	if err := r.billing.Run(ctx.Context()); err != nil {
		r.l.Error("run billing check:", zap.Error(err))
	}
}

// getVacantProperties godoc
// @Summary      Search vacant properties
// @Description  Protected. Returns every property system-wide that currently has no lease bound to it, along with the applications submitted for each.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Success      200  {array}   entity.PropertyDetail
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/vacant [get]
func (r *V1) getVacantProperties(ctx fiber.Ctx) error {
	properties, err := r.p.GetVacantProperties(ctx.Context())
	if err != nil {
		r.l.Error("get vacant properties:", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get vacant properties")
	}

	return ctx.Status(http.StatusOK).JSON(properties)
}

// applyToProperty godoc
// @Summary      Apply for a property
// @Description  Protected, Tenant only. Submits the caller's application to rent the given property.
// @Tags         properties
// @Produce      json
// @Security     CookieAuth
// @Param        id   path      string  true  "Property ID"
// @Success      200  {object}  map[string]bool
// @Failure      400  {object}  response.Error  "already applied to this property"
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      404  {object}  response.Error  "property not found"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/{id}/apply [post]
func (r *V1) applyToProperty(ctx fiber.Ctx) error {
	tenantUserID, _ := ctx.Locals(middleware.UserIDLocalsKey).(string)

	if err := r.p.Apply(ctx.Context(), ctx.Params("id"), tenantUserID); err != nil {
		switch {
		case errors.Is(err, repo.ErrPropertyNotFound):
			return errorResponse(ctx, http.StatusNotFound, "property not found")
		case errors.Is(err, repo.ErrApplicationAlreadyExists):
			return errorResponse(ctx, http.StatusBadRequest, "Вы уже откликнулись на это предложение")
		default:
			r.l.Error("apply to property:", zap.Error(err))
			return errorResponse(ctx, http.StatusInternalServerError, "failed to submit application")
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}
