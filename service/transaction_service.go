package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/web"
	"github.com/google/uuid"
)

type TransactionServiceImpl struct {
	db *sql.DB
	q  *repository.Queries
}

func NewTransactionService(dbCon *sql.DB) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		db: dbCon,
		q:  repository.New(dbCon),
	}
}

func (service *TransactionServiceImpl) CreateTransaction(ctx context.Context, arg repository.CreateTransactionParams) (web.TransactionResponse, error) {
	var resp web.TransactionResponse

	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to begin transaction: %w", err)
	}

	txQueries := repository.New(tx)
	req := repository.GetLimitParams{
		CustomerID: arg.CustomerID,
		Tenor:      arg.Tenor,
	}

	limit, err := txQueries.GetLimit(ctx, req)
	if err != nil {
		tx.Rollback()
		return resp, exception.NewNotFoundError("Doesn't have limit")
	}

	if limit < arg.Price {
		tx.Rollback()
		return resp, exception.NewNotFoundError(fmt.Sprintf("Limit tidak cukup: Max %.2f", limit))
	}

	if err := txQueries.CreateTransaction(ctx, arg); err != nil {
		tx.Rollback()
		return resp, exception.NewNotFoundError(err.Error())
	}

	totalPrice := arg.Price + arg.Bunga + arg.AdminFee
	paymentAmount := totalPrice / float64(arg.Tenor)

	for i := 1; i <= int(arg.Tenor); i++ {
		err := txQueries.CreatePayment(ctx, repository.CreatePaymentParams{
			ID:            uuid.NewString(),
			TransactionID: arg.ID,
			Amount:        paymentAmount,
			DueDate:       arg.CreatedAt.AddDate(0, i, 0),
			IsPaid:        false,
		})
		if err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			return resp, fmt.Errorf("failed to insert payment detail for month %d: %w", i+1, err)
		}
	}

	err = txQueries.ReduceLimit(ctx, repository.ReduceLimitParams{
		Limit:      totalPrice,
		CustomerID: arg.CustomerID,
		Limit_2:    totalPrice,
	})

	if err != nil {
		tx.Rollback()
		return resp, fmt.Errorf("failed to reduce limit: %w", err)
	}

	payloadResp, err := txQueries.GetTransaction(ctx, arg.ID)
	if err != nil {
		tx.Rollback()
		return resp, exception.NewNotFoundError(err.Error())
	}

	if err := tx.Commit(); err != nil {
		return resp, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return web.NewTransactionResponse(payloadResp), nil
}

func (service *TransactionServiceImpl) ListTx(ctx context.Context, id string) ([]web.TransactionResponse, error) {
	var resp []web.TransactionResponse
	payload, err := service.q.ListTransaction(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, fmt.Errorf("no data founded : %s", err.Error())
		}
		return resp, exception.NewNotFoundError(err.Error())
	}
	return web.NewTransactionResponses(payload), nil
}
