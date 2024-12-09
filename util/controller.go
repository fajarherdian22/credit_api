package util

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func ValidateDate(input string) (time.Time, error) {
	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, input)

	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, must be YYYY-MM-DD")
	}
	return parsedDate, nil
}
