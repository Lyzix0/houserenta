package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/potom_pridumaem/internal/controller/v1/response"
	"go.uber.org/zap"
)

// getUnlinkedTenants godoc
// @Summary      List unlinked tenants
// @Description  Protected, Landlord only. Returns registered tenants not currently bound to any property, for picking one when manually creating a lease.
// @Tags         tenants
// @Produce      json
// @Security     CookieAuth
// @Success      200  {array}   response.TenantSummary
// @Failure      401  {object}  response.Error  "not authenticated"
// @Failure      500  {object}  response.Error  "internal server error"
// @Router       /tenants/unlinked [get]
func (r *V1) getUnlinkedTenants(ctx fiber.Ctx) error {
	tenants, err := r.u.GetUnlinkedTenants(ctx.Context())
	if err != nil {
		r.l.Error("get unlinked tenants:", zap.Error(err))
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get unlinked tenants")
	}

	summaries := make([]response.TenantSummary, 0, len(tenants))
	for _, tenant := range tenants {
		summaries = append(summaries, response.TenantSummary{
			ID:       tenant.ID,
			Name:     tenant.Name,
			Email:    tenant.Email,
			Document: tenant.Document,
			Phone:    tenant.Phone,
		})
	}

	return ctx.Status(http.StatusOK).JSON(summaries)
}
