package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	"github.com/potom_pridumaem/internal/repo"
	"go.uber.org/zap"
)

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
