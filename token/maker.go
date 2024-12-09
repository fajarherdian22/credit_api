package token

import "time"

type Maker interface {
	CreateToken(email, customerId string, duration time.Duration) (string, *Payload, error)

	VerifiyToken(token string) (*Payload, error)
}
