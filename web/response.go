package web

import (
	"time"

	"github.com/fajarherdian22/credit_bank/repository"
)

type CustomerResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func NewCustomerResponse(customer repository.Customer) CustomerResponse {
	return CustomerResponse{
		ID:       customer.ID,
		FullName: customer.FullName,
		Email:    customer.Email,
	}
}

type LoginUserResponse struct {
	SessionID             string           `json:"session_id"`
	AccessToken           string           `json:"access_token"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	Customer              CustomerResponse `json:"customer"`
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
