package service

import (
	"context"
	"database/sql"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type LoanServiceImpl struct {
	q *repository.Queries
}

func NewLoanService(q *repository.Queries) *LoanServiceImpl {
	return &LoanServiceImpl{
		q: q,
	}
}

func (service *LoanServiceImpl) CreateLimit(ctx context.Context, customerID string) error {
	salary, err := service.q.GetSalary(ctx, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return exception.NewBadRequestError("salary not found for customer")
		}
		return exception.NewInternalError(err.Error())
	}
	originalSalary := salary
	for i := 1; i <= 4; i++ {
		salary = originalSalary
		switch i {
		case 1:
			salary *= 0.15
		case 2:
			salary *= 0.2
		case 3:
			salary *= 0.25
		case 4:
			salary *= 0.3
		}
		arg := repository.GenerateLimitParams{
			ID:         uuid.NewString(),
			CustomerID: customerID,
			Tenor:      int32(i),
			Limit:      salary,
		}

		err := service.q.GenerateLimit(ctx, arg)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				switch mysqlErr.Number {
				case 1062:
					return exception.NewForbiddenError("duplicate limit entry")
				case 1451:
					return exception.NewBadRequestError("foreign key constraint violation")
				}
			}
			return exception.NewInternalError(err.Error())
		}
	}
	return nil
}

func (service *LoanServiceImpl) ListLimit(ctx context.Context, customer_id string) ([]repository.GetCustomerLimitRow, error) {
	var resp []repository.GetCustomerLimitRow
	payload, err := service.q.GetCustomerLimit(ctx, customer_id)
	if err != nil {
		return resp, exception.NewNotFoundError(err.Error())
	}
	return payload, nil
}
