package api

import (
	"context"
	"konnex/pkg/rest"
	"konnex/users"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type Endpoint struct {
	RegisterEndpoint  endpoint.Endpoint
	AuthorizeEndpoint endpoint.Endpoint
}

func MakeServerEndpoint(svc users.Service) Endpoint {
	return Endpoint{
		RegisterEndpoint:  RegisterEndpoint(svc),
		AuthorizeEndpoint: AuthorizeEndpoint(svc),
	}
}

func RegisterEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserReqBody)

		NewUser := &users.User{
			Username: req.Username,
			Password: req.Password,
		}

		err = svc.Register(ctx, *NewUser)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusInternalServerError,
				Status:  "Error",
				Message: "Err Create User",
				Errors:  err.Error(),
			}, nil
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

		CurrentUser := &users.User{
			Username: req.Username,
			Password: req.Password,
		}

		auth, err := svc.Authorize(ctx, *CurrentUser)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusBadRequest,
				Status:  "Error",
				Message: "Err Authorize",
				Errors:  err.Error(),
			}, nil
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
