package controller

import (
	"net/http"
	"time"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/model/web"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionController struct {
	TransactionService *service.TransactionServiceImpl
}

func NewTransactionController(TransactionService *service.TransactionServiceImpl) *TransactionController {
	return &TransactionController{TransactionService: TransactionService}
}

type CreateTransactionsRequest struct {
	CustomerID  string  `json:"customer_id" binding:"required,len=36"`
	ProductName string  `json:"product_name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Tenor       int32   `json:"tenor" binding:"required,gt=0"`
}

func createTransactionsPayload(req CreateTransactionsRequest, total service.TotalPayment) repository.CreateTransactionParams {
	return repository.CreateTransactionParams{
		ID:            uuid.NewString(),
		CustomerID:    req.CustomerID,
		ProductName:   req.ProductName,
		Price:         req.Price,
		Bunga:         total.Bunga,
		JumlahCicilan: total.JumlahCicilan,
		Tenor:         req.Tenor,
		AdminFee:      total.AdminFee,
		CreatedAt:     time.Now(),
	}
}

func (controller *TransactionController) CreateTransaction(c *gin.Context) {
	var req CreateTransactionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	total := service.CalculateTotalPayment(req.Price, req.Tenor)

	arg := createTransactionsPayload(req, total)

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
