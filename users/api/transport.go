package api

import (
	"context"
	"encoding/json"
	"errors"
	"konnex/pkg/rest"
	"konnex/users"
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

func MakeHTTPHandler(svc users.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoint(svc)
	opt := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r.Methods("POST").Path("/register").Handler(httptransport.NewServer(
		e.RegisterEndpoint,
		decodeUserRequest,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("POST").Path("/authorize").Handler(httptransport.NewServer(
		e.AuthorizeEndpoint,
		decodeUserRequest,
		rest.EncodeResponse,
		opt...,
	))

	return r
}

func decodeUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req UserReqBody
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

type UserReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
