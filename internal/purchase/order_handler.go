package purchase

import (
	"belimang/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	usecase IOrderUsecase
}

func NewOrderHandler(uc IOrderUsecase) *orderHandler {
	return &orderHandler{
		usecase: uc,
	}
}

func (h *orderHandler) Router(r *gin.RouterGroup) {
	// Order Route Group
	group := r.Group(
		"users",
		// middleware.UseJwtAuth,
	)

	// Routing
	group.POST("estimate", h.Estimate)
}

func (h *orderHandler) Estimate(c *gin.Context) {
	var req Request

	// Parse request body to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	// Validate the starting point
	if err := req.ValidateRequest(); err != nil {
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	// Estimate user order via usecase
	result, err := h.usecase.Estimate(req)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	response.GenerateResponse(c, 200, response.WithData(result))
}
