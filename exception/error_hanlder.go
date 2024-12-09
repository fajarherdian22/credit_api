package exception

import (
	"net/http"

	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/model/web"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

type AppError struct {
	Code    int
	Status  string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
func NewAppError(code int, status, message string) *AppError {
	return &AppError{Code: code, Status: status, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, "NOT FOUND", message)
}

func NewNotAuthError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, "BAD REQUEST", message)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, "FORBIDDEN", message)
}

func NewInternalError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERRROR", message)
}

func ErrorHandler(c *gin.Context, err interface{}) {

	if handleAppError(c, err) {
		return
	}

	if handleValidationErrors(c, err) {
		return
	}
}

func handleAppError(c *gin.Context, err interface{}) bool {
	appError, ok := err.(*AppError)
	if ok {
		sendErrorResponse(c, appError.Code, appError.Status, appError.Message)
		return true
	}
	return false
}

func handleValidationErrors(c *gin.Context, err interface{}) bool {
	validationErrs, ok := err.(validator.ValidationErrors)
	if ok {
		sendErrorResponse(c, http.StatusBadRequest, "VALIDATION ERROR", validationErrs.Error())
		return true
	}
	return false
}

func sendErrorResponse(c *gin.Context, code int, status string, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)

	webResponse := web.WebResponse{
		Code:   code,
		Status: status,
		Data:   data,
	}

	helper.HandleEncodeWriteJson(c, webResponse)
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
