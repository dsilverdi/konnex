package users

import (
	"context"
	"time"
)

type Auth struct {
	ID          string
	AccessToken string
	Expiration  int
	CreatedAt   time.Time
}

type AuthRepository interface {
	Save(ctx context.Context, token, id string) error
	Authorize(ctx context.Context, id string) (*Auth, error)
}
