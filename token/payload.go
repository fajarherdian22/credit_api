package token

import (
	"errors"
	"time"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrInvalid = errors.New("Token Invalid")
	ErrExpired = errors.New("Token Expired")
)

type Payload struct {
	ID         string    `json:id`
	Email      string    `json:email`
	CustomerID string    `json:customer_id`
	IssuedAt   time.Time `json:issued_at`
	ExpiredAt  time.Time `json:expired_at`
}

func NewPayload(email, customerId string, duration time.Duration) (*Payload, error) {

	payload := &Payload{
		ID:         uuid.NewString(),
		Email:      email,
		CustomerID: customerId,
		IssuedAt:   time.Now(),
		ExpiredAt:  time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpired
	}
	return nil
}

func GetPayload(c *gin.Context) (*Payload, error) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		return nil, exception.NewNotAuthError("not authenticated")
	}

	tokenPayload, ok := payload.(*Payload)
	if !ok {
		return nil, exception.NewInternalError("failed to parse token payload")
	}

	return tokenPayload, nil
}
