package middleware

import (
	"belimang/pkg/jwt"
	"belimang/pkg/response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Protect routes using JWT scheme.
// Get Authorization header from request, then validate the token using JWT.
// Token has Bearer type.
func UseJwtAuth(ctx *gin.Context) {
	// Get token
	// If token is not valid, return error
	token, err := getAuthorizationToken(ctx)
	if err != nil {
		response.GenerateResponse(ctx, http.StatusUnauthorized, response.WithMessage(err.Error()))
		ctx.Abort()

		return
	}

	// Validate Token
	// If token cannot validated, return error
	result, err := jwt.ValidateToken(token)
	if err != nil {
		response.GenerateResponse(ctx, http.StatusUnauthorized, response.WithMessage(err.Error()))
		ctx.Abort()

		return
	}

	// After token successfully validated,
	// set UserID from token to current context
	ctx.Set("userID", result.Uuid)
	ctx.Set("userRole", result.Role)

	// Next if passed middleware
	ctx.Next()
}

// Get Bearer token from request header.
// Return error if no token provided or scheme is not Bearer
func getAuthorizationToken(ctx *gin.Context) (string, error) {
	// Get Authorizatin from header
	authHeader := ctx.GetHeader("Authorization")

	// Check if token is in bearer format
	// Split string ideally return slice with 2 element, the 2nd element is the token
	bearerToken := strings.Split(authHeader, "Bearer ")
	if len(bearerToken) == 2 {
		return bearerToken[1], nil
	}

	return "", errors.New("token is not valid")
}
