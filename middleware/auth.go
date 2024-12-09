package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationBearer     = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("error authorization header is not provide")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid auth format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}
		authType := strings.ToLower(fields[0])
		if authType != authorizationBearer {
			err := fmt.Errorf("unsupported auth type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifiyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.ErrorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
