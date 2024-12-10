package web

import (
	"time"

	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/google/uuid"
)

type CreateCustomersRequest struct {
	Nik          string  `json:"nik" binding:"required,len=16"`
	Password     string  `json:"password" binding:"required,min=6"`
	Email        string  `json:"email" binding:"required,email"`
	FullName     string  `json:"full_name" binding:"required"`
	LegalName    string  `json:"legal_name" binding:"required"`
	TempatLahir  string  `json:"tempat_lahir" binding:"required"`
	TanggalLahir string  `json:"tanggal_lahir" binding:"required"`
	Gaji         float64 `json:"gaji" binding:"required,numeric,gt=0"`
	PhotoKtp     string  `json:"photo_ktp" binding:"required"`
	FotoSelfie   string  `json:"foto_selfie" binding:"required"`
}

func CreateCustomersPayload(req CreateCustomersRequest, pw string, tgl_lahir time.Time) repository.CreateCustomersParams {
	return repository.CreateCustomersParams{
		ID:             uuid.NewString(),
		Nik:            req.Nik,
		HashedPassword: pw,
		Email:          req.Email,
		FullName:       req.FullName,
		LegalName:      req.LegalName,
		TempatLahir:    req.TempatLahir,
		TanggalLahir:   tgl_lahir,
		Gaji:           req.Gaji,
		PhotoKtp:       req.PhotoKtp,
		FotoSelfie:     req.FotoSelfie,
	}
}

type CreateTransactionsRequest struct {
	CustomerID  string  `json:"customer_id" binding:"required,len=36"`
	ProductName string  `json:"product_name" binding:"required,ProductName"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Tenor       int32   `json:"tenor" binding:"required,gt=0"`
}

type TotalPayment struct {
	Bunga         float64
	JumlahCicilan float64
	AdminFee      float64
}

func CreateTransactionsPayload(req CreateTransactionsRequest, total TotalPayment) repository.CreateTransactionParams {
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
