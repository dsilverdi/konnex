// Copyright (c) Konnex

package users

import (
	"context"
	"time"
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

type UserJSONResponse struct {
	ID        string
	Username  string
	CreatedAt time.Time
}

type UserRepository interface {
	Save(ctx context.Context, user User) error
	Read(ctx context.Context, username string) (*User, error)
	ReadbyID(ctx context.Context, id string) (*User, error)
}
