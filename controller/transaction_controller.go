package controller

import (
	"net/http"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/token"
	"github.com/fajarherdian22/credit_bank/web"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TransactionController struct {
	TransactionService *service.TransactionServiceImpl
	tokenMaker         token.Maker
	Validate           *validator.Validate
}

func NewTransactionController(TransactionService *service.TransactionServiceImpl, tokenMaker token.Maker, validate *validator.Validate) *TransactionController {
	return &TransactionController{
		TransactionService: TransactionService,
		tokenMaker:         tokenMaker,
		Validate:           validate,
	}
}

func getPayload(c *gin.Context) *token.Payload {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		exception.ErrorHandler(c, exception.NewNotAuthError("not authenticated"))
	}
	tokenPayload, ok := payload.(*token.Payload)
	if !ok {
		exception.ErrorHandler(c, exception.NewInternalError("failed to parse token payload"))
	}
	return tokenPayload
}

func CalculateTotalPayment(price float64, tenor int32) web.TotalPayment {
	bunga := 0.1
	total := price + (price * bunga)
	jumlahCicilan := total / float64(tenor)
	adminFee := jumlahCicilan * 0.15

	return web.TotalPayment{
		Bunga:         bunga,
		JumlahCicilan: jumlahCicilan,
		AdminFee:      adminFee,
	}
}

func (controller *TransactionController) CreateTransaction(c *gin.Context) {
	var req web.CreateTransactionsRequest
	tokenPayload := getPayload(c)
	req.CustomerID = tokenPayload.CustomerID

	if err := c.ShouldBindJSON(&req); err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	if err := controller.Validate.Struct(req); err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	total := CalculateTotalPayment(req.Price, req.Tenor)

	arg := web.CreateTransactionsPayload(req, total)

	payload, err := controller.TransactionService.CreateTransaction(c.Request.Context(), arg)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	WebResponse := web.WebResponse{
		Code:   http.StatusOK,
		Data:   payload,
		Status: "OK",
	}

	helper.HandleEncodeWriteJson(c, WebResponse)
}

func (controller *TransactionController) ListTx(c *gin.Context) {
	tokenPayload := getPayload(c)
	payload, err := controller.TransactionService.ListTx(c, tokenPayload.CustomerID)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	WebResponse := web.WebResponse{
		Code:   http.StatusOK,
		Data:   payload,
		Status: "OK",
	}

	helper.HandleEncodeWriteJson(c, WebResponse)
}
