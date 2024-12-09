package service

import (
	"time"

	"github.com/fajarherdian22/credit_bank/repository"
)

type TotalPayment struct {
	Bunga         float64
	JumlahCicilan float64
	AdminFee      float64
}

func CalculateTotalPayment(price float64, tenor int32) TotalPayment {
	bunga := 0.1
	total := price + (price * bunga)
	jumlahCicilan := total / float64(tenor)
	adminFee := jumlahCicilan * 0.15

	return TotalPayment{
		Bunga:         bunga,
		JumlahCicilan: jumlahCicilan,
		AdminFee:      adminFee,
	}
}

type TransactionResponse struct {
	ID                string    `json:"transaction_id"`
	CustomerID        string    `json:"customer_id"`
	ProductName       string    `json:"product_name"`
	TotalPrice        float64   `json:"total_price"`
	TotalInstallments float64   `json:"total_installments"`
	Tenor             int32     `json:"tenor"`
	Interest          float64   `json:"interest"`
	AdminFee          float64   `json:"admin_fee"`
	CreatedAt         time.Time `json:"transaction_at"`
}

func NewTransactionResponse(customers repository.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:                customers.ID,
		CustomerID:        customers.CustomerID,
		ProductName:       customers.ProductName,
		TotalPrice:        customers.Price,
		TotalInstallments: customers.JumlahCicilan,
		Tenor:             customers.Tenor,
		Interest:          customers.Bunga,
		AdminFee:          customers.AdminFee,
		CreatedAt:         customers.CreatedAt,
	}
}
func NewTransactionResponses(transaction []repository.Transaction) []TransactionResponse {
	var txResponses []TransactionResponse
	for _, tx := range transaction {
		txResponses = append(txResponses, NewTransactionResponse(tx))
	}
	return txResponses
}
