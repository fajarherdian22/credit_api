package helper

import "github.com/gin-gonic/gin"

func IsError(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
