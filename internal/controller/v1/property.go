package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/potom_pridumaem/internal/controller/middleware"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo/persistent"
)

// getProperties godoc
// @Summary      List user properties
// @Description  Returns properties owned by the authenticated landlord
// @Tags         properties
// @Produce      json
// @Success      200  {array}   entity.Property
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties [get]
func (r *V1) getProperties(ctx fiber.Ctx) error {
	userID := ctx.Locals(middleware.UserIDLocalsKey).(string)
	props, err := r.p.GetProperties(ctx.Context(), userID)
	if err != nil {
		r.l.Error("restapi - v1 - getProperties", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	if props == nil {
		props = []entity.Property{}
	}
	return ctx.Status(http.StatusOK).JSON(props)
}

// createProperty godoc
// @Summary      Create a property
// @Description  Creates a new rental property for a landlord
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
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	prop, err := r.p.CreateProperty(ctx.Context(), body)
	if err != nil {
		r.l.Error("restapi - v1 - createProperty", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"id": prop.ID, "name": prop.Name})
}

// getProperty godoc
// @Summary      Get property by ID
// @Description  Returns a single property by its ID
// @Tags         properties
// @Produce      json
// @Param        id   path      string  true  "Property ID"
// @Success      200  {object}  entity.Property
// @Failure      404  {object}  response.Error  "property not found"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/{id} [get]
func (r *V1) getProperty(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	prop, err := r.p.GetProperty(ctx.Context(), id)
	if err != nil {
		r.l.Error("restapi - v1 - getProperty", zap.Error(err))
		return errorResponse(ctx, http.StatusNotFound, "property not found")
	}
	return ctx.Status(http.StatusOK).JSON(prop)
}

// updateProperty godoc
// @Summary      Update a property
// @Description  Updates an existing property
// @Tags         properties
// @Accept       json
// @Produce      json
// @Param        id    path      string           true  "Property ID"
// @Param        input  body      request.Property  true  "Property data"
// @Success      200    {object}  map[string]bool
// @Failure      400    {object}  response.Error  "invalid request body"
// @Failure      500    {object}  response.Error  "internal server error"
// @Router       /properties/{id} [put]
func (r *V1) updateProperty(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var body request.Property
	if err := ctx.Bind().Body(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := r.p.UpdateProperty(ctx.Context(), id, body); err != nil {
		r.l.Error("restapi - v1 - updateProperty", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// deleteProperty godoc
// @Summary      Delete a property
// @Description  Deletes a property by ID
// @Tags         properties
// @Produce      json
// @Param        id   path      string  true  "Property ID"
// @Success      200  {object}  map[string]bool
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/{id} [delete]
func (r *V1) deleteProperty(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := r.p.DeleteProperty(ctx.Context(), id); err != nil {
		r.l.Error("restapi - v1 - deleteProperty", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"ok": true})
}

// getVacantProperties godoc
// @Summary      List vacant properties
// @Description  Returns properties with no active lease
// @Tags         properties
// @Produce      json
// @Success      200  {array}   entity.Property
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /properties/vacant [get]
func (r *V1) getVacantProperties(ctx fiber.Ctx) error {
	props, err := r.p.GetVacantProperties(ctx.Context())
	if err != nil {
		r.l.Error("restapi - v1 - getVacantProperties", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	if props == nil {
		props = []entity.Property{}
	}
	return ctx.Status(http.StatusOK).JSON(props)
}

func (r *V1) getUnlinkedTenants(ctx fiber.Ctx) error {
	tenants, err := r.propertyRepo.GetUnlinkedTenants(ctx.Context())
	if err != nil {
		r.l.Error("restapi - v1 - getUnlinkedTenants", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	if tenants == nil {
		tenants = []persistent.UnlinkedTenant{}
	}
	return ctx.Status(http.StatusOK).JSON(tenants)
}
