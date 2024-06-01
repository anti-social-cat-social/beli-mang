package purchase

import (
	"belimang/internal/middleware"
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
		middleware.UseJwtAuth,
	)

	// Routing
	group.POST("estimate", h.Estimate)
	group.POST("orders", h.Order)
	group.GET("orders", h.OrderHistory)
}

func (h *orderHandler) Estimate(c *gin.Context) {
	var req Request

	userId := c.GetString("userID")

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
	req.UserId = userId

	result, err := h.usecase.Estimate(req)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	response.GenerateResponse(c, 200, response.WithData(result))
}

func (h *orderHandler) Order(c *gin.Context) {
	var entity ActualOrder

	// Parse request body to struct
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	result, err := h.usecase.PlaceOrder(entity)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	response.GenerateResponse(c, 201, response.WithData(result))
}

func (h *orderHandler) OrderHistory(c *gin.Context) {
	var req GetOrderHistQueryParams

	userId := c.GetString("userID")

	// Parse request body to struct
	if err := c.ShouldBindQuery(&req); err != nil {
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	result, err := h.usecase.OrderHistory(userId, req)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	response.GenerateResponse(c, 200, response.WithData(result))
}