package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseOpts func(*response) error

// Default success response struct
type response struct {
	Code    int         `json:"-"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

// Customize response message
// by default it is using http status StatusText
// ex : code 201 - message Created
func WithMessage(message interface{}) responseOpts {
	return func(r *response) error {
		r.Message = message

		return nil
	}
}

// Add data to the response
// By default it won't return any data, data = nil
// Data can be single value, map, or struct
func WithData(data any) responseOpts {
	return func(r *response) error {
		r.Data = data

		return nil
	}
}

// Return JSON response to the caller
// code should be defined. options is used to customize the JSON response.
func GenerateResponse(ctx *gin.Context, code int, options ...responseOpts) {
	response := response{
		Code:    code,
		Data:    nil,
		Message: http.StatusText(code),
	}

	for _, opts := range options {
		opts(&response)
	}

	ctx.JSON(
		code,
		response,
	)
}
