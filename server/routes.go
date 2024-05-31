package server

import (
	"belimang/internal/merchant"
	"belimang/internal/purchase"
	"belimang/internal/user"
	"belimang/internal/image"
	"belimang/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRoute(engine *gin.Engine, db *sqlx.DB) {
	// Handle for not found routes
	engine.NoRoute(NoRouteHandler)
	router := engine.Group("v1")

	router.GET("ping", pingHandler)

	initializeMerchantHandler(db, router)
	initializeUserHandler(db, router)
	initializeOrderHandler(db, router)
	initializeImageHandler(router)
}

func initializeMerchantHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initialize all necessary dependecies
	merchantRepo := merchant.NewMerchantRepository(db)
	merchantUc := merchant.NewMerchantUsecase(merchantRepo)
	merchantH := merchant.NewMerchantHandler(merchantUc)

	merchantH.Router(router)
}

func initializeUserHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initialize all necessary dependecies
	userRepo := user.NewUserRepository(db)
	userUc := user.NewUserUsecase(userRepo)
	userH := user.NewUserHandler(userUc)

	userH.Router(router)
}

func initializeOrderHandler(db *sqlx.DB, router *gin.RouterGroup) {
	merchantRepo := merchant.NewMerchantRepository(db)
	merchantUc := merchant.NewMerchantUsecase(merchantRepo)

	orderRepo := purchase.NewOrderRepository(db)
	orderUc := purchase.NewOrderUsecase(orderRepo, merchantUc)
	orderH := purchase.NewOrderHandler(orderUc)

	orderH.Router(router)
}

func initializeImageHandler(router *gin.RouterGroup) {
	imageH := image.NewImageHandler()

	imageH.Router(router)
}

func NoRouteHandler(ctx *gin.Context) {
	response.GenerateResponse(ctx, http.StatusNotFound, response.WithMessage("Page not found"))
}

// Handler for ping request from routes
func pingHandler(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		struct {
			Data    any    `json:"data"`
			Message string `json:"message"`
			Success bool   `json:"success"`
		}{
			Success: true,
			Message: "Server is online",
			Data:    true,
		},
	)
}
