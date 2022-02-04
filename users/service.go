package users

import (
	"context"
	"crypto/sha256"
	"fmt"
	"konnex"
	"konnex/pkg/errors"
	"time"
)

//JSON Format Struct

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

var (
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrCreateUUID indicates error in creating uuid for entity creation
	ErrCreateUUID = errors.New("uuid creation failed")

	// ErrCreateEntity indicates error in creating entity or entities
	ErrCreateEntity = errors.New("create entity failed")

	// ErrUpdateEntity indicates error in updating entity or entities
	ErrUpdateEntity = errors.New("update entity failed")

	// ErrAuthorization indicates a failure occurred while authorizing the entity.
	ErrAuthorization = errors.New("failed to perform authorization over the entity")

	// ErrViewEntity indicates error in viewing entity or entities
	ErrViewEntity = errors.New("view entity failed")

	// ErrRemoveEntity indicates error in removing entity
	ErrRemoveEntity = errors.New("remove entity failed")

	// ErrConnect indicates error in adding connection
	ErrConnect = errors.New("add connection failed")

	// ErrDisconnect indicates error in removing connection
	ErrDisconnect = errors.New("remove connection failed")

	// ErrFailedToRetrieveThings failed to retrieve things.
	ErrFailedToRetrieveThings = errors.New("failed to retrieve group members")

	// ErrWrongPassword indicates error in wrong password
	ErrWrongPassword = errors.New("Wrong Password")
)

type Service interface {
	// Register New User
	Register(context.Context, User) error

	// Login User
	Authorize(context.Context, User) (*Auth, error)

	// Get Authorized User Data
	ViewAccount(context.Context, string) (*User, error)
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
	NewUser := &User{
		Username:  user.Username,
		Password:  fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))),
		CreatedAt: time.Now(),
	}

	id, err := svc.IDProvider.ID()
	if err != nil {
		return errors.Wrap(ErrCreateUUID, err)
	}

	NewUser.ID = id

	// Perform DB Call Here
	err = svc.UserRepo.Save(ctx, *NewUser)
	if err != nil {
		return errors.Wrap(ErrCreateEntity, err)
	}

	return nil
}

func (svc *UserService) Authorize(ctx context.Context, user User) (*Auth, error) {
	CurrentUser, err := svc.UserRepo.Read(ctx, user.Username)
	if err != nil {
		return nil, errors.Wrap(ErrNotFound, err)
	}

	password := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	if password != CurrentUser.Password {
		return nil, errors.Wrap(ErrWrongPassword, err)
	}

	// Perform DB Call for Auth Repo Here
	auth, err := svc.AuthRepo.Authorize(ctx, CurrentUser.ID)
	if err != nil {
		return nil, errors.Wrap(ErrViewEntity, err)
	}

	if auth == nil {
		NewToken, err := svc.IDProvider.ID()
		if err != nil {
			return nil, errors.Wrap(ErrCreateUUID, err)
		}

		auth = &Auth{
			ID:          CurrentUser.ID,
			AccessToken: NewToken,
			Expiration:  0,
			CreatedAt:   time.Now(),
		}
	}

	return auth, nil
}

func (svc *UserService) ViewAccount(ctx context.Context, token string) (*User, error) {
	return nil, nil
}
