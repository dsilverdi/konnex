package grpc

import (
	"context"
	"konnex"
	"konnex/pkg/errors"
	"konnex/users"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	authorize kitgrpc.Handler
	identify  kitgrpc.Handler
}

func NewServer(svc users.Service) konnex.AuthServiceServer {
	return &grpcServer{
		authorize: kitgrpc.NewServer(
			authorizeEndpoint(svc),
			decodeAuthorizeRequest,
			encodeAuthorizeResponse,
		),
		identify: kitgrpc.NewServer(
			identifyEndpoint(svc),
			decodeIdentifyRequest,
			encodeIdentifyResponse,
		),
	}
}

func (s *grpcServer) Authorize(ctx context.Context, token *konnex.Token) (*konnex.AuthorizeRes, error) {
	_, res, err := s.authorize.ServeGRPC(ctx, token)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*konnex.AuthorizeRes), nil
}

func (s *grpcServer) Identify(ctx context.Context, id *konnex.UserID) (*konnex.UserIdentity, error) {
	_, res, err := s.identify.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*konnex.UserIdentity), nil
}

func encodeError(err error) error {
	switch {
	case errors.Contains(err, nil):
		return nil
	case errors.Contains(err, errors.ErrMalformedEntity):
		return status.Error(codes.InvalidArgument, "received invalid token request")
	case errors.Contains(err, errors.ErrUnauthorizedAccess),
		errors.Contains(err, errors.ErrAuthorization):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
