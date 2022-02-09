package grpc

import (
	"context"
	"konnex"
	"konnex/pkg/errors"
)

type authorizeReq struct {
	token string
}

func (req authorizeReq) validate() error {
	if req.token == "" {
		return errors.ErrMalformedEntity
	}

	return nil
}

func decodeAuthorizeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*konnex.Token)
	return authorizeReq{
		token: req.Value,
	}, nil
}

func encodeAuthorizeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(authorizeReq)
	return &konnex.Token{
		Value: req.token,
	}, nil
}

type identifyReq struct {
	userid string
}

func (req identifyReq) validate() error {
	if req.userid == "" {
		return errors.ErrMalformedEntity
	}

	return nil
}

func decodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*konnex.UserID)
	return identifyReq{
		userid: req.Value,
	}, nil
}

func encodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(identifyReq)
	return &konnex.UserID{
		Value: req.userid,
	}, nil
}
