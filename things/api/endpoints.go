package api

import (
	"context"
	"fmt"
	"konnex/things"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateThingEndpoint endpoint.Endpoint
	GetThingsEndpoint   endpoint.Endpoint

	CreateChannelEndpoint endpoint.Endpoint
	GetChannelEndpoint    endpoint.Endpoint
}

func MakeServerEndpoint(svc things.Service) Endpoints {
	return Endpoints{
		CreateThingEndpoint: CreateThingsEndpoint(svc),
		GetThingsEndpoint:   GetThingsEndpoint(svc),

		CreateChannelEndpoint: CreateChannelEndpoint(svc),
		GetChannelEndpoint:    GetChannelEndpoint(svc),
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
		if err != nil {
			return nil, err
		}

		res := createThingsResponse{
			Things:  *th,
			Message: "Success",
			Err:     e,
		}
		return res, nil
	}
}

func GetThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getThingsReq)
		var things []things.Things

		fmt.Println(req.channelID)

		things, err = svc.GetThings(ctx)
		if err != nil {
			return nil, err
		}

		resp := getThingsRes{
			Things:  things,
			Message: "Success",
		}
		return resp, nil
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
			return nil, err
		}

		resp := createChannelResponse{
			Channel: *ch,
			Message: "success",
			Err:     err,
		}

		return resp, nil
	}
}

func GetChannelEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getChannelReq)
		var channels []things.Channel

		fmt.Println(req.Type)

		channels, err = svc.GetChannels(ctx)
		if err != nil {
			return nil, err
		}

		resp := getChannelResponse{
			Message:  "success",
			Channels: channels,
		}

		return resp, nil
	}
}
