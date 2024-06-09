package user

import (
	// "belimang/internal/middleware"
	// "belimang/pkg/jwt"
	"belimang/pkg/response"
	"belimang/pkg/validation"
	"net/http"
	// "log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	uc IUserUsecase
}

// Constructor for user handler struct
func NewUserHandler(uc IUserUsecase) *userHandler {
	return &userHandler{
		uc: uc,
	}
}

func (h *userHandler) Router(r *gin.RouterGroup) {
	// Grouping to give URL prefix
	userRoute := r.Group("users")
	adminRoute := r.Group("admin")

	// route for users
	userRoute.POST("login", h.Login(USER))
	userRoute.POST("register", h.Register(USER))

	// route for admin
	adminRoute.POST("login", h.Login(ADMIN))
	adminRoute.POST("register", h.Register(ADMIN))

	// group.GET("", middleware.UseJwtAuth, middleware.HasRoles(string(IT)), h.GetUsers)
}

func (h *userHandler) Login(r UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UserLoginDTO
		if err := ctx.ShouldBindJSON(&request); err != nil {
			response.GenerateResponse(ctx, 400)
			ctx.Abort()
			return
		}

		// Validate request
		validate := validator.New(validator.WithRequiredStructEnabled())

		// Generate error validation if not any field is not valid
		if err := validate.Struct(request); err != nil {
			validatorMessage := validation.GenerateStructValidationError(err)
			response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
			ctx.Abort()
			return
		}

		requestData := UserLoginWithRoleDTO{
			Username: request.Username,
			Password: request.Password,
			Role: string(r),
		}

		resp, respError := h.uc.Login(requestData)
		if respError != nil {
			response.GenerateResponse(ctx, respError.Code, response.WithMessage(respError.Error.Error()))
			ctx.Abort()
			return
		}

		response.GenerateResponseReturnData(ctx, 200, response.WithData(*resp))
	}
}

func (h *userHandler) Register(r UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UserRegisterDTO
		if err := ctx.ShouldBindJSON(&request); err != nil {
			// validatorMessage := validation.GenerateStructValidationError(err)
			// response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
			response.GenerateResponse(ctx, 400)
			ctx.Abort()
			return
		}

		// Validate request
		validate := validator.New(validator.WithRequiredStructEnabled())

		// Generate error validation if not any field is not valid
		if err := validate.Struct(request); err != nil {
			validatorMessage := validation.GenerateStructValidationError(err)
			response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
			ctx.Abort()
			return
		}

		requestData := UserRegisterWithRoleDTO{
			Username: request.Username,
			Password: request.Password,
			Email: request.Email,
			Role: string(r),
		}

		resp, respError := h.uc.Register(requestData)
		if respError != nil {
			response.GenerateResponse(ctx, respError.Code, response.WithMessage(respError.Error.Error()))
			ctx.Abort()
			return
		}

		response.GenerateResponseReturnData(ctx, 201, response.WithData(*resp))
	}
}