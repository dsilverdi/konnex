package grpc

import (
	"context"
	"konnex"
)

type authorizeRes struct {
	token  string
	userid string
}

func decodeAuthorizeResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*konnex.AuthorizeRes)
	return authorizeRes{
		token:  res.Token,
		userid: res.UserID,
	}, nil
}

func encodeAuthorizeResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(authorizeRes)
	return konnex.AuthorizeRes{
		Token:  res.token,
		UserID: res.userid,
	}, nil
}

type identityRes struct {
	userid   string
	username string
}

func decodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*konnex.UserIdentity)
	return identityRes{
		userid:   res.Id,
		username: res.Username,
	}, nil
}

func encodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(identityRes)
	return konnex.UserIdentity{
		Id:       res.userid,
		Username: res.username,
	}, nil
}
