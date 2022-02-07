package users

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"konnex"
	"konnex/pkg/errors"
	"time"

	stderr "errors"

	"github.com/go-sql-driver/mysql"
)

//JSON Format Struct

type Service interface {
	// Register New User
	Register(context.Context, User) error

	// Login User
	Authorize(context.Context, User) (*Auth, error)

	// Get Authorized User Data
	ViewAccount(context.Context, string) (*User, error)

	// Token Validation
	TokenValidation(context.Context, string) (*string, error)
}

type UserService struct {
	UserRepo   UserRepository
	AuthRepo   AuthRepository
	IDProvider konnex.IDprovider
}

func New(userrepo UserRepository, authrepo AuthRepository, idprovider konnex.IDprovider) Service {
	return &UserService{
		UserRepo:   userrepo,
		AuthRepo:   authrepo,
		IDProvider: idprovider,
	}
}

func (svc *UserService) Register(ctx context.Context, user User) error {
	var mysqlErr *mysql.MySQLError

	NewUser := &User{
		Username:  user.Username,
		Password:  fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))),
		CreatedAt: time.Now(),
	}

	id, err := svc.IDProvider.ID()
	if err != nil {
		return errors.Wrap(errors.ErrCreateUUID, err)
	}

	NewUser.ID = id

	// Perform DB Call Here
	err = svc.UserRepo.Save(ctx, *NewUser)
	if err != nil {
		if stderr.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return errors.Wrap(errors.ErrAlreadyExists, err)
		}
		return errors.Wrap(errors.ErrCreateEntity, err)
	}

	return nil
}

func (svc *UserService) Authorize(ctx context.Context, user User) (*Auth, error) {
	CurrentUser, err := svc.UserRepo.Read(ctx, user.Username)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNotFound, err)
	}

	password := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	if password != CurrentUser.Password {
		return nil, errors.ErrWrongPassword
	}

	// Perform DB Call for Auth Repo Here
	auth, err := svc.AuthRepo.Authorize(ctx, CurrentUser.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			NewToken, err := svc.IDProvider.ID()
			if err != nil {
				return nil, errors.Wrap(errors.ErrCreateUUID, err)
			}

			auth = &Auth{
				ID:          CurrentUser.ID,
				AccessToken: NewToken,
				Expiration:  0,
				CreatedAt:   time.Now(),
			}

			err = svc.AuthRepo.Save(ctx, *auth)
			if err != nil {
				return nil, errors.Wrap(errors.ErrCreateEntity, err)
			}
		} else {
			return nil, errors.Wrap(errors.ErrViewEntity, err)
		}
	}

	return auth, nil
}

func (svc *UserService) ViewAccount(ctx context.Context, token string) (*User, error) {
	id, err := svc.TokenValidation(ctx, token)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := svc.UserRepo.ReadbyID(ctx, *id)
	if err != nil {
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return user, nil
}

func (svc *UserService) TokenValidation(ctx context.Context, token string) (*string, error) {
	id, err := svc.AuthRepo.Validate(ctx, token)
	if err != nil {
		return nil, err
	}

	return id, nil
}
