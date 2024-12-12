package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/go-sql-driver/mysql"
)

type CustomerServiceImpl struct {
	q *repository.Queries
}

func NewCustomerService(q *repository.Queries) *CustomerServiceImpl {
	return &CustomerServiceImpl{
		q: q,
	}
}

type CreateCustomersResponse struct {
	FullName     string    `json:"full_name"`
	LegalName    string    `json:"legal_name"`
	Nik          string    `json:"nik"`
	TempatLahir  string    `json:"tempat_lahir"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	Email        string    `json:"email"`
	Created_at   time.Time `json:"created_at"`
}

func NewCustomerPayload(customers repository.Customer) CreateCustomersResponse {
	return CreateCustomersResponse{
		FullName:     customers.FullName,
		LegalName:    customers.LegalName,
		Nik:          customers.Nik,
		TempatLahir:  customers.TempatLahir,
		TanggalLahir: customers.TanggalLahir,
		Email:        customers.Email,
		Created_at:   time.Now(),
	}
}

func (service *CustomerServiceImpl) GetCustomer(ctx context.Context, arg string) (repository.Customer, error) {
	var resp repository.Customer
	payload, err := service.q.GetCustomers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, exception.NewBadRequestError(err.Error())
		}
		return resp, exception.NewInternalError(err.Error())
	}
	return payload, nil
}

func (service *CustomerServiceImpl) CreateCustomers(ctx context.Context, arg repository.CreateCustomersParams) (CreateCustomersResponse, error) {
	var resp CreateCustomersResponse
	err := service.q.CreateCustomers(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				return resp, exception.NewForbiddenError("Duplicate entry error")
			case 1451:
				return resp, exception.NewBadRequestError(mysqlErr.Message)
			}
		}
		return resp, exception.NewNotFoundError(err.Error())
	}

	payload, err := service.q.GetCustomers(ctx, arg.Email)
	if err != nil {
		return resp, exception.NewNotFoundError(err.Error())
	}

	return NewCustomerPayload(payload), nil
}

func (service *CustomerServiceImpl) CreateSession(ctx context.Context, arg repository.CreateSessionParams) (repository.Session, error) {
	var resp repository.Session
	err := service.q.CreateSession(ctx, arg)
	if err != nil {
		return resp, err
	}

	payload, err := service.q.GetSession(ctx, arg.ID)
	if err != nil {
		return resp, err
	}
	return payload, nil
}

func (service *CustomerServiceImpl) GetSession(ctx context.Context, id string) (repository.Session, error) {
	var resp repository.Session
	payload, err := service.q.GetSession(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, exception.NewNotFoundError(err.Error())
		}
		return resp, err
	}
	return payload, nil
}
