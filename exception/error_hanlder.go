package exception

import (
	"net/http"

	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/model/web"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func ErrorHandler(c *gin.Context, err interface{}) {

	if notFoundError(c, err) {
		return
	}

	if validationErrors(c, err) {
		return
	}

	internalServerError(c, err)
}

func validationErrors(c *gin.Context, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Data:   exception.Error(),
			Status: "Bad Request",
		}

		helper.HandleEncodeWriteJson(c, webResponse)
		return true
	} else {
		return false
	}
}

func notFoundError(c *gin.Context, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error,
		}

		helper.HandleEncodeWriteJson(c, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(c *gin.Context, err interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err,
	}

	helper.HandleEncodeWriteJson(c, webResponse)
}
