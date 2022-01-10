package api

import (
	"context"
	"fmt"
	"konnex/pkg/rest"
	"konnex/things"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateThingEndpoint       endpoint.Endpoint
	GetThingsEndpoint         endpoint.Endpoint
	GetSpecificThingsEndpoint endpoint.Endpoint
	DeleteThingEndpoint       endpoint.Endpoint

	CreateChannelEndpoint      endpoint.Endpoint
	GetChannelEndpoint         endpoint.Endpoint
	GetSpecificChannelEndpoint endpoint.Endpoint
	DeleteChannelEndpoint      endpoint.Endpoint
}

func MakeServerEndpoint(svc things.Service) Endpoints {
	return Endpoints{
		CreateThingEndpoint:       CreateThingsEndpoint(svc),
		GetThingsEndpoint:         GetThingsEndpoint(svc),
		GetSpecificThingsEndpoint: GetSpecificThingsEndpoint(svc),
		DeleteThingEndpoint:       DeleteThingEndpoint(svc),

		CreateChannelEndpoint:      CreateChannelEndpoint(svc),
		GetChannelEndpoint:         GetChannelEndpoint(svc),
		GetSpecificChannelEndpoint: GetSpecificChannelEndpoint(svc),
		DeleteChannelEndpoint:      DeleteChannelEndpoint(svc),
	}
}

func CreateThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createThingsReq)

		data := things.Things{
			ID:        req.ID,
			ChannelID: req.ChannelID,
			Name:      req.Name,
			MetaData:  req.MetaData,
		}

		th, e := svc.CreateThings(ctx, data)
		if e != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Create Things Error",
				Errors:  e.Error(),
			}, e
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Create Things Success",
			Data:    th,
		}, nil
	}
}

func GetThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getThingsReq)
		var things []things.Things

		fmt.Println(req.channelID)

		things, err = svc.GetThings(ctx)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Get Things Error",
				Errors:  err.Error(),
			}, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Get List of Things",
			Data:    things,
		}, nil
	}
}

func GetSpecificThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getSpecificReq)

		things, err := svc.GetSpecificThing(ctx, req.ID)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Get Things Error",
				Errors:  err.Error(),
			}, nil
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Get Specific Thing",
			Data:    things,
		}, nil
	}
}

func DeleteThingEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getSpecificReq)

		if req.ID == "" {
			return rest.HTTPResponse{
				Code:    http.StatusBadRequest,
				Status:  "No ID provided",
				Message: "Delete Things Error",
			}, nil
		}

		err = svc.DeleteThing(ctx, req.ID)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Delete Things Error",
				Errors:  err.Error(),
			}, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Success Delete Thing",
		}, nil
	}
}

// Channel Endpoint

func CreateChannelEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createChannelReq)

		data := things.Channel{
			ID:       req.ID,
			Name:     req.Name,
			Type:     req.Type,
			Metadata: req.MetaData,
		}

		ch, err := svc.CreateChannel(ctx, data)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Create Channel Error",
				Errors:  err.Error(),
			}, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Create Channel Success",
			Data:    ch,
		}, nil
	}
}

func GetChannelEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getChannelReq)
		var channels []things.Channel

		fmt.Println(req.Type)

		channels, err = svc.GetChannels(ctx)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Get Channel Error",
				Errors:  err.Error(),
			}, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Get Channel Success",
			Data:    channels,
		}, nil
	}
}

func GetSpecificChannelEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getSpecificReq)

		ch, err := svc.GetSpecificChannel(ctx, req.ID)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Create Channel Error",
				Errors:  err.Error(),
			}, nil
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Get Channel Success",
			Data:    ch,
		}, nil
	}
}

func DeleteChannelEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getSpecificReq)

		if req.ID == "" {
			return rest.HTTPResponse{
				Code:    http.StatusBadRequest,
				Status:  "No ID provided",
				Message: "Delete Channel Error",
			}, nil
		}

		err = svc.DeleteChannel(ctx, req.ID)
		if err != nil {
			return rest.HTTPResponse{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: "Delete Channel Error",
				Errors:  err.Error(),
			}, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Success Delete Channel",
		}, nil
	}
}
