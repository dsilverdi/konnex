package data

import (
	"context"
	"konnex/users"
	"time"
)

type UserRepository struct {
	db Database
}
type UserDB struct {
	ID        string    `db:"id"`
	UserName  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func NewChannelRepository(db Database) users.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (db *UserRepository) Save(ctx context.Context, user users.User) error {
	return nil
}

func (db *UserRepository) Read(ctx context.Context, username string) (*users.User, error) {
	return nil, nil
}
