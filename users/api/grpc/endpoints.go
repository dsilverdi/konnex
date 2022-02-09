package grpc

import (
	"context"
	"konnex/users"

	"github.com/go-kit/kit/endpoint"
)

func identifyEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(identifyReq)

		if err := req.validate(); err != nil {
			return identityRes{}, nil
		}

		user, err := svc.IdentifyUser(ctx, req.userid)
		if err != nil {
			return identityRes{}, err
		}

		res := identityRes{
			userid:   user.ID,
			username: user.Username,
		}

		return res, nil
	}
}

func authorizeEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authorizeReq)

		if err := req.validate(); err != nil {
			return authorizeRes{}, err
		}

		id, err := svc.TokenValidation(ctx, req.token)
		if err != nil {
			return authorizeRes{}, err
		}

		res := authorizeRes{
			userid: *id,
			token:  req.token,
		}

		return res, nil
	}
}
