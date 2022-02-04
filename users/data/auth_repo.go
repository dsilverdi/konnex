package data

import (
	"context"
	"konnex/users"
	"time"
)

type AuthRepository struct {
	db Database
}
type UserAuthDB struct {
	ID          string    `db:"id"`
	AccessToken string    `db:"access_token"`
	Expired     string    `db:"expired"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewAuthRepository(db Database) users.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

// create token for new login
func (db *AuthRepository) Save(ctx context.Context, token, id string) error {
	return nil
}

// get token for login
func (db *AuthRepository) Authorize(ctx context.Context, id string) (*users.Auth, error) {
	return nil, nil
}
