package token

import (
	"errors"
	"time"

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
