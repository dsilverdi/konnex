package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"konnex/opcua"
	"konnex/pkg/rest"
	"net/http"

	"github.com/go-kit/log"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler ")
)

func MakeHTTPHandler(svc opcua.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoint(svc)
	opt := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(rest.EncodeError),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r.Methods("POST").Path("/browse").Handler(httptransport.NewServer(
		e.BrowseEndpoint,
		decodeBrowseRequest,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("GET").Path("/monitor/{id}").Handler(httptransport.NewServer(
		e.MonitorEndpoint,
		decodeGetData,
		rest.EncodeResponse,
		opt...,
	))
	fmt.Print(opt)
	return r
}

func decodeBrowseRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req BrowseReq
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetData(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	// token := r.Header.Get("Authorization")

	return GetMonitorReq{
		ID: id,
	}, nil
}

type BrowseReq struct {
	ServerURI  string `json:"server-uri"`
	NameSpace  string `json:"namespace"`
	Identifier string `json:"identifier"`
}

type GetMonitorReq struct {
	ID string
}
