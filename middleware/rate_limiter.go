package middleware

import (
	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(2, 5)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			exception.ErrorHandler(c, exception.NewManyRequest("Too many requests"))
			c.Abort()
			return
		}
		c.Next()
	}
}
