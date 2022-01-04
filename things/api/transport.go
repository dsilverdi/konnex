package api

import (
	"context"
	"encoding/json"
	"konnex/things"
	"net/http"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(svc things.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoint(svc)
	opt := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	// POST    /things/                          adds things
	// GET     /things/:id                       retrieves the given things by id
	// PUT     /things/:id                       post updated things information about the things
	// PATCH   /things/:id                       partial updated things information
	// DELETE  /things/:id                       remove the given things
	// GET     /things/:id/addresses/            retrieve addresses associated with the things
	// GET     /things/:id/addresses/:addressID  retrieve a particular things address
	// POST    /things/:id/addresses/            add a new address
	// DELETE  /things/:id/addresses/:addressID  remove an address

	r.Methods("POST").Path("/things/").Handler(httptransport.NewServer(
		e.AddThingsEndpoint,
		decodeAddThingsRequest,
		encodeResponse,
		opt...,
	))

	r.Methods("GET").Path("/things/").Handler(httptransport.NewServer(
		e.GetThingsEndpoint,
		decodeGetThingsRequest,
		encodeResponse,
		opt...,
	))

	return r
}

func decodeAddThingsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req things.PostThingsRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Things); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetThingsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req things.PostThingsRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Things); e != nil {
		return nil, e
	}
	return req, nil
}

type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
// func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
// 	var buf bytes.Buffer
// 	err := json.NewEncoder(&buf).Encode(request)
// 	if err != nil {
// 		return err
// 	}
// 	req.Body = ioutil.NopCloser(&buf)
// 	return nil
// }

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case things.ErrNotFound:
		return http.StatusNotFound
	case things.ErrAlreadyExists, things.ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
