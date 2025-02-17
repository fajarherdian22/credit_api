// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"context"
)

type Querier interface {
	CreateCustomers(ctx context.Context, arg CreateCustomersParams) error
	CreatePayment(ctx context.Context, arg CreatePaymentParams) error
	CreateSession(ctx context.Context, arg CreateSessionParams) error
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) error
	GenerateLimit(ctx context.Context, arg GenerateLimitParams) error
	GetCustomerLimit(ctx context.Context, customerID string) ([]GetCustomerLimitRow, error)
	GetCustomers(ctx context.Context, email string) (Customer, error)
	GetLimit(ctx context.Context, arg GetLimitParams) (float64, error)
	GetSalary(ctx context.Context, id string) (float64, error)
	GetSession(ctx context.Context, id string) (Session, error)
	GetTransaction(ctx context.Context, id string) (Transaction, error)
	ListTransaction(ctx context.Context, customerID string) ([]Transaction, error)
	ReduceLimit(ctx context.Context, arg ReduceLimitParams) error
}

var _ Querier = (*Queries)(nil)
