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
	Expired     int       `db:"expired"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewAuthRepository(db Database) users.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

// create token for new login
func (a *AuthRepository) Save(ctx context.Context, auth users.Auth) error {
	authDB := &UserAuthDB{
		ID:          auth.ID,
		AccessToken: auth.AccessToken,
		Expired:     auth.Expiration,
		CreatedAt:   auth.CreatedAt,
	}

	query := `INSERT INTO users_auth (id, access_token, expired, created_at)
	VALUES (:id, :access_token, :expired, :created_at);`

	_, err := a.db.NamedExecContext(ctx, query, authDB)
	if err != nil {
		return err
	}

	return nil
}

// get token for login
func (a *AuthRepository) Authorize(ctx context.Context, id string) (*users.Auth, error) {
	var authDB UserAuthDB

	query := `SELECT id, access_token, expired, created_at FROM users_auth WHERE id = ?`

	err := a.db.QueryRowxContext(ctx, query, id).StructScan(&authDB)
	if err != nil {
		return nil, err
	}

	Auth := &users.Auth{
		ID:          authDB.ID,
		AccessToken: authDB.AccessToken,
		Expiration:  authDB.Expired,
		CreatedAt:   authDB.CreatedAt,
	}

	return Auth, nil
}

// validate token
func (a *AuthRepository) Validate(ctx context.Context, token string) (*string, error) {
	var authDB UserAuthDB

	query := `SELECT id, access_token, expired, created_at FROM users_auth WHERE access_token = ?`

	err := a.db.QueryRowxContext(ctx, query, token).StructScan(&authDB)
	if err != nil {
		return nil, err
	}

	return &authDB.ID, nil
}
