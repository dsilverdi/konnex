package grpc

import (
	"context"
	"konnex"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

const (
	svcName = "konnex.AuthService"
)

type grpcClient struct {
	authorize endpoint.Endpoint
	identify  endpoint.Endpoint
	timeout   time.Duration
}

func NewClient(conn *grpc.ClientConn, timeout time.Duration) konnex.AuthServiceClient {
	return &grpcClient{
		authorize: kitgrpc.NewClient(
			conn,
			svcName,
			"Authorize",
			encodeAuthorizeRequest,
			decodeAuthorizeResponse,
			konnex.AuthorizeRes{},
		).Endpoint(),

		identify: kitgrpc.NewClient(
			conn,
			svcName,
			"Identify",
			encodeIdentifyRequest,
			decodeIdentifyResponse,
			konnex.UserIdentity{},
		).Endpoint(),

		timeout: timeout,
	}
}

func (cl *grpcClient) Authorize(ctx context.Context, req *konnex.Token, _ ...grpc.CallOption) (*konnex.AuthorizeRes, error) {
	ctx, close := context.WithTimeout(ctx, cl.timeout)
	defer close()

	res, err := cl.authorize(ctx, authorizeReq{token: req.GetValue()})
	if err != nil {
		return nil, err
	}

	ar := res.(authorizeRes)
	return &konnex.AuthorizeRes{UserID: ar.userid, Token: ar.token}, nil
}

func (cl *grpcClient) Identify(ctx context.Context, req *konnex.UserID, _ ...grpc.CallOption) (*konnex.UserIdentity, error) {
	ctx, close := context.WithTimeout(ctx, cl.timeout)
	defer close()

	res, err := cl.identify(ctx, identifyReq{userid: req.GetValue()})
	if err != nil {
		return nil, err
	}

	ir := res.(identityRes)
	return &konnex.UserIdentity{Id: ir.userid, Username: ir.username}, nil
}
