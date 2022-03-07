package rest

import (
	"context"
	"konnex/pkg/errors"
	"konnex/pkg/rest"
	"konnex/users"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type Endpoint struct {
	RegisterEndpoint   endpoint.Endpoint
	AuthorizeEndpoint  endpoint.Endpoint
	GetProfileEndpoint endpoint.Endpoint
}

func MakeServerEndpoint(svc users.Service) Endpoint {
	return Endpoint{
		RegisterEndpoint:   RegisterEndpoint(svc),
		AuthorizeEndpoint:  AuthorizeEndpoint(svc),
		GetProfileEndpoint: GetProfileEndpoint(svc),
	}
}

func RegisterEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserReqBody)

		if req.Username == "" {
			return nil, errors.ErrCreateEntity
		}

		if req.Password == "" {
			return nil, errors.ErrCreateEntity
		}

		NewUser := &users.User{
			Username: req.Username,
			Password: req.Password,
		}

		err = svc.Register(ctx, *NewUser)
		if err != nil {
			return nil, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Success Register User",
		}, nil
	}
}

func AuthorizeEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserReqBody)

		if req.Username == "" {
			return nil, errors.ErrCreateEntity
		}

		if req.Password == "" {
			return nil, errors.ErrCreateEntity
		}

		CurrentUser := &users.User{
			Username: req.Username,
			Password: req.Password,
		}

		auth, err := svc.Authorize(ctx, *CurrentUser)
		if err != nil {
			return nil, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Success Authorize User",
			Data: map[string]string{
				"token": auth.AccessToken,
			},
		}, nil

	}
}

func GetProfileEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetProfileReq)

		if req.Token == "" {
			return nil, errors.ErrUnauthorizedAccess
		}

		profile, err := svc.ViewAccount(ctx, req.Token)
		if err != nil {
			return nil, err
		}

		return rest.HTTPResponse{
			Code:   http.StatusOK,
			Status: "Success",
			Data:   profile,
		}, nil
	}
}
