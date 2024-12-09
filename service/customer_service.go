package service

import (
	"context"
	"time"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/repository"
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
	payload, err := service.q.GetCustomers(ctx, arg)
	if err != nil {
		return repository.Customer{}, err
	}

	return payload, nil
}

func (service *CustomerServiceImpl) CreateCustomers(ctx context.Context, arg repository.CreateCustomersParams) (CreateCustomersResponse, error) {
	err := service.q.CreateCustomers(ctx, arg)
	if err != nil {
		return CreateCustomersResponse{}, exception.NewNotFoundError(err.Error())
	}

	payload, err := service.q.GetCustomers(ctx, arg.ID)
	if err != nil {
		return CreateCustomersResponse{}, exception.NewNotFoundError(err.Error())
	}

	return NewCustomerPayload(payload), nil
}

func (service *CustomerServiceImpl) CreateSession(ctx context.Context, arg repository.CreateSessionParams) (repository.Session, error) {
	err := service.q.CreateSession(ctx, arg)
	if err != nil {
		return repository.Session{}, err
	}

	payload, err := service.q.GetSession(ctx, arg.ID)
	if err != nil {
		return repository.Session{}, err
	}
	return payload, nil

}
