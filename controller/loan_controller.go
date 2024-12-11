package controller

import (
	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/token"
	"github.com/fajarherdian22/credit_bank/web"
	"github.com/gin-gonic/gin"
)

type LoanController struct {
	LoanService *service.LoanServiceImpl
	tokenMaker  token.Maker
}

func NewLoanController(LoanService *service.LoanServiceImpl, tokenMaker token.Maker) *LoanController {
	return &LoanController{
		LoanService: LoanService,
		tokenMaker:  tokenMaker,
	}
}

func (controller *LoanController) GenerateLimit(c *gin.Context) {
	tokenPayload, err := token.GetPayload(c)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	err = controller.LoanService.CreateLimit(c, tokenPayload.CustomerID)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	payload, err := controller.LoanService.ListLimit(c, tokenPayload.CustomerID)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	limits := web.NewLimitResponses(payload)

	response := web.GenerateLimitResponse{
		CustomerID: tokenPayload.CustomerID,
		Email:      tokenPayload.Email,
		Limits:     limits,
	}

	WebResponse := web.WebResponse{
		Code:   200,
		Data:   response,
		Status: "OK",
	}

	helper.HandleEncodeWriteJson(c, WebResponse)
}
