package localError

import (
	"halosuster/pkg/logger"
	"fmt"
	"net/http"
	"runtime"
)

type GlobalError struct {
	Code    int
	Message interface{}
	Error   error
}

// Return Not found error structure with customize message and error.
func ErrBase(code int, message interface{}, err error) *GlobalError {
	baseError := GlobalError{
		Code:    code,
		Message: message,
		Error:   err,
	}

	return &baseError
}

// Return internal server error structure with customize message and error.
func ErrInternalServer(message string, err error) *GlobalError {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		logger.Error(err, fmt.Sprintf("%s:%d", filename, line))
	} else {
		logger.Info(message)
	}

	baseError := ErrBase(http.StatusInternalServerError, message, err)

	return baseError
}

// Return unauthorized structure with customize message and error.
func ErrUnauthorized(message string, err error) *GlobalError {
	if err != nil {
		logger.Info(err.Error())
	} else {
		logger.Info(message)
	}

	baseError := ErrBase(http.StatusUnauthorized, message, err)

	return baseError
}

// Return unauthorized structure with customize message and error.
func ErrForbidden(message string, err error) *GlobalError {
	if err != nil {
		logger.Info(err.Error())
	} else {
		logger.Info(message)
	}

	baseError := ErrBase(http.StatusForbidden, message, err)

	return baseError
}

// Return Not found error structure with customize message and error.
func ErrNotFound(message string, err error) *GlobalError {
	if err != nil {
		logger.Info(err.Error())
	} else {
		logger.Info(message)
	}

	baseError := ErrBase(http.StatusNotFound, message, err)

	return baseError
}

// Return conflict error structure with customize message and error.
func ErrConflict(message string, err error) *GlobalError {
	if err != nil {
		logger.Info(err.Error())
	} else {
		logger.Info(message)
	}

	baseError := ErrBase(http.StatusConflict, message, err)

	return baseError
}

// Return bad request error structure with customize message and error.
func ErrBadRequest(message interface{}, err error) *GlobalError {
	if err != nil {
		logger.Info(err.Error())
	} else {
		logger.Info(message)
	}

	baseError := ErrBase(http.StatusBadRequest, message, err)

	return baseError
}
