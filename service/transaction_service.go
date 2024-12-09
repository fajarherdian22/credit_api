package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/repository"
)

type TransactionServiceImpl struct {
	q *repository.Queries
}

func NewTransactionService(q *repository.Queries) *TransactionServiceImpl {
	return &TransactionServiceImpl{q: q}
}

func (service *TransactionServiceImpl) CreateTransaction(ctx context.Context, arg repository.CreateTransactionParams) (TransactionResponse, error) {
	var resp TransactionResponse
	req := repository.GetLimitParams{
		CustomerID: arg.CustomerID,
		Tenor:      arg.Tenor,
	}

	limit, err := service.q.GetLimit(ctx, req)
	if err != nil {
		return resp, exception.NewNotFoundError(err.Error())
	}

	if limit < arg.Price {
		return resp, exception.NewNotFoundError(fmt.Sprintf("Limit tidak cukup: Max %.2f", limit))
	}

	if err := service.q.CreateTransaction(ctx, arg); err != nil {
		return resp, exception.NewNotFoundError(err.Error())
	}
	payloadResp, err := service.q.GetTransaction(ctx, arg.ID)
	if err != nil {
		return resp, exception.NewNotFoundError(err.Error())
	}
	return NewTransactionResponse(payloadResp), nil
}

func (service *TransactionServiceImpl) ListTx(ctx context.Context, id string) ([]TransactionResponse, error) {
	var resp []TransactionResponse
	payload, err := service.q.ListTransaction(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, fmt.Errorf("no data founded : %s", err.Error())
		}
		return resp, exception.NewNotFoundError(err.Error())
	}
	return NewTransactionResponses(payload), nil
}
