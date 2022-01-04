package api

import (
	"context"
	"konnex/things"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddThingsEndpoint endpoint.Endpoint
	GetThingsEndpoint endpoint.Endpoint
}

func MakeServerEndpoint(svc things.Service) Endpoints {
	return Endpoints{
		AddThingsEndpoint: MakeAddThingsEndpoint(svc),
		GetThingsEndpoint: MakeGetThingsEndpoint(svc),
	}
}

func MakeAddThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(things.PostThingsRequest)
		e := svc.AddThings(ctx, req.Things)
		return things.PostThingsResponse{Err: e}, nil
	}
}

func MakeGetThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// req := request.(things.PostThingsRequest)
		e := svc.GetThings(ctx)
		return things.PostThingsResponse{Err: e}, nil
	}
}
