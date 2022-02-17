package api

import (
	"context"
	"konnex/opcua"
	"konnex/pkg/rest"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type Endpoint struct {
	BrowseEndpoint  endpoint.Endpoint
	MonitorEndpoint endpoint.Endpoint
}

func MakeServerEndpoint(svc opcua.Service) Endpoint {
	return Endpoint{
		BrowseEndpoint:  BrowseEndpoint(svc),
		MonitorEndpoint: MonitorEndpoint(svc),
	}
}

func BrowseEndpoint(svc opcua.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(BrowseReq)

		nodes, err := svc.Browse(ctx, req.ServerURI, req.NameSpace, req.Identifier)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Browse Node Error",
				Errors:  err.Error(),
			}, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Browse Node List",
			Data:    nodes,
		}, nil
	}
}

func MonitorEndpoint(svc opcua.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetMonitorReq)
		datas, err := svc.Monitor(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return rest.HTTPResponse{
			Code:   http.StatusOK,
			Status: "Success",
			Total:  len(datas),
			Data:   datas,
		}, nil
	}
}
