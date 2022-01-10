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
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r.Methods("GET").Path("/browse").Handler(httptransport.NewServer(
		e.BrowseEndpoint,
		decodeBrowseRequest,
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

type BrowseReq struct {
	ServerURI  string `json:"server-uri"`
	NameSpace  string `json:"namespace"`
	Identifier string `json:"identifier"`
}
