package merchant

import (
	"belimang/internal/middleware"
	"belimang/internal/user"
	"belimang/pkg/jwt"
	"belimang/pkg/response"
	"belimang/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type merchantHandler struct {
	uc IMerchantUsecase
}

// Constructor for user handler struct
func NewMerchantHandler(uc IMerchantUsecase) *merchantHandler {
	return &merchantHandler{
		uc: uc,
	}
}

func (h *merchantHandler) Router(r *gin.RouterGroup) {
	// Grouping to give URL prefix
	group := r.Group("admin/merchants")

	group.POST("", middleware.UseJwtAuth, middleware.HasRoles(string(user.ADMIN)), h.CreateMerchant)
	group.POST("/:merchantId/items", middleware.UseJwtAuth, middleware.HasRoles(string(user.ADMIN)), h.CreateItem)
}

func (h *merchantHandler) CreateMerchant(ctx *gin.Context) {
	var request CreateMerchantDTO

	if err := ctx.ShouldBindJSON(&request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		return
	}

	resp, respError := h.uc.CreateMerchant(request)
	if respError != nil {
		response.GenerateResponse(ctx, respError.Code, response.WithMessage(respError.Error.Error()))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithData(*resp))
}

func (h *merchantHandler) CreateItem(ctx *gin.Context) {
	var request CreateItemDTO
	merchantId := ctx.Param("merchantId")

	if err := ctx.ShouldBindJSON(&request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		return
	}

	resp, respError := h.uc.CreateItem(merchantId, request)
	if respError != nil {
		response.GenerateResponse(ctx, respError.Code, response.WithMessage(respError.Error.Error()))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithData(*resp))
}