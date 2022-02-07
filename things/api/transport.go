package api

import (
	"context"
	"encoding/json"
	"errors"
	"konnex/pkg/rest"
	"konnex/things"
	"net/http"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler ")
)

func MakeHTTPHandler(svc things.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoint(svc)
	opt := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(rest.EncodeError),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	// POST    	/things/                         create things
	// GET	   	/things/						 get list of things
	// GET     	/things/:id                      retrieves the given things by id
	// DELETE  	/things/:id                      remove the given things

	// POST		/channel/						 create channel
	// GET		/channel/						 get list of channel
	// GET		/channel/:id					 get specific channel
	// DELETE	/channel/:id					 delete channel

	// Things Functionality Route
	// Things Define the Device / Sensor Data / IoT Node that client want to Observe

	r.Methods("GET").Path("/things").Handler(httptransport.NewServer(
		e.GetThingsEndpoint,
		decodeGetThingsRequest,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("POST").Path("/things/").Handler(httptransport.NewServer(
		e.CreateThingEndpoint,
		decodeCreateThingsRequest,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("GET").Path("/things/{id}").Handler(httptransport.NewServer(
		e.GetSpecificThingsEndpoint,
		decodeGetSpecific,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("DELETE").Path("/things/{id}").Handler(httptransport.NewServer(
		e.DeleteThingEndpoint,
		decodeGetSpecific,
		rest.EncodeResponse,
		opt...,
	))

	// Channel Functionality Route
	// Channel Define Group of Connectivity that client want to Observe
	// e.g MQTT, OPCUA, HTTP, etc..

	r.Methods("GET").Path("/channel").Handler(httptransport.NewServer(
		e.GetChannelEndpoint,
		decodeGetChannelEndpoint,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("GET").Path("/channel/{id}").Handler(httptransport.NewServer(
		e.GetSpecificChannelEndpoint,
		decodeGetSpecific,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("POST").Path("/channel/").Handler(httptransport.NewServer(
		e.CreateChannelEndpoint,
		decodeCreateChannelRequest,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("DELETE").Path("/channel/{id}").Handler(httptransport.NewServer(
		e.DeleteChannelEndpoint,
		decodeGetSpecific,
		rest.EncodeResponse,
		opt...,
	))

	// Connection Functionality Route

	// r.Methods("PATCH").Path("/connect/{id}").Handler(httptransport.NewServer(

	// ))

	// r.Methods("PATCH").Path("/disconnect/{id}").Handler(httptransport.NewServer(

	// ))

	return r
}

func decodeCreateThingsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createThingsReq
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetThingsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	channelID := r.URL.Query().Get("channel_id")
	return getThingsReq{channelID: channelID}, nil
}

func decodeCreateChannelRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createChannelReq

	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetChannelEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {

	Type := r.URL.Query().Get("type")

	return getChannelReq{Type: Type}, nil
}

func decodeGetSpecific(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getSpecificReq{ID: id}, nil
}
