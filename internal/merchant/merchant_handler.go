package merchant

import (
	"belimang/internal/middleware"
	"belimang/internal/user"
	"belimang/pkg/response"
	"belimang/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	"log"
	"strings"
	"strconv"
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
	adminGroup := r.Group("admin/merchants")
	userGroup := r.Group("")

	adminGroup.POST("", middleware.UseJwtAuth, middleware.HasRoles(string(user.ADMIN)), h.CreateMerchant)
	adminGroup.POST("/:merchantId/items", middleware.UseJwtAuth, middleware.HasRoles(string(user.ADMIN)), h.CreateItem)
	adminGroup.GET("", middleware.UseJwtAuth, middleware.HasRoles(string(user.ADMIN)), h.FindAllMerchants)

	userGroup.GET("/merchants/nearby/:latlong", middleware.UseJwtAuth, middleware.HasRoles(string(user.USER)), h.GetLatLong, h.FindNearbyMerchants)
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

func (h *merchantHandler) FindAllMerchants(c *gin.Context) {
	query := GetMerchantQueryParams{}
	if err := c.ShouldBindQuery(&query); err != nil {
		res := validation.FormatValidation(err)
		response.GenerateResponse(c, res.Code, response.WithMessage(res.Message))
		return
	}

	merchants, err := h.uc.FindAllMerchants(query)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("Product fetched successfully!"), response.WithData(merchants))
}

func (h *merchantHandler) GetLatLong(ctx *gin.Context) {
	latlong := ctx.Param("latlong")
	latlongArr := strings.Split(latlong, ",")

	if len(latlongArr) != 2 {
		response.GenerateResponse(ctx, 400, response.WithMessage("lat / long is not valid"))
		ctx.Abort()
		return
	}

	lat, err := strconv.ParseFloat(latlongArr[0], 32)
	if err != nil {
		response.GenerateResponse(ctx, 400, response.WithMessage("lat / long is not valid"))
		ctx.Abort()
		return
	}

	long, err := strconv.ParseFloat(latlongArr[1], 32)
	if err != nil {
		response.GenerateResponse(ctx, 400, response.WithMessage("lat / long is not valid"))
		ctx.Abort()
		return
	}

	ctx.Set("location", Location{
		Lat: float32(lat),
		Long: float32(long),
	})

	ctx.Next()
}

func (h *merchantHandler) FindNearbyMerchants(c *gin.Context) {
	query := GetMerchantQueryParams{}	

	if err := c.ShouldBindQuery(&query); err != nil {
		res := validation.FormatValidation(err)
		response.GenerateResponse(c, res.Code, response.WithMessage(res.Message))
		return
	}

	var location Location
	locationI, _ := c.Get("location")
	location = locationI.(Location)
	// log.Println(location)

	// merchants := Merchant{}
	// merchants, err := h.uc.FindAllMerchants(query)
	// if err != nil {
	// 	response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
	// 	return
	// }

	merchants, err := h.uc.FindNearbyMerchants(location, query)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		return
	}

	log.Println(merchants)

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("Product fetched successfully!"), response.WithData(merchants))
}
