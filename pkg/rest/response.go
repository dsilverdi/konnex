package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"konnex/pkg/errors"
	"net/http"
)

const (
	contentType = "application/json"
)

type HTTPResponse struct {
	Code    int         `json:"code,omitempty"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Total   int         `json:"total,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type errorer interface {
	error() error
}

type errorRes struct {
	Err string `json:"error"`
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		EncodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.

func EncodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errorVal := err.(type) {
	case errors.Error:
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		switch {
		case errors.Contains(errorVal, errors.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusUnauthorized)

		case errors.Contains(errorVal, errors.ErrAuthorization):
			w.WriteHeader(http.StatusForbidden)
		// case errors.Contains(errorVal, errors.ErrInvalidQueryParams):
		// 	w.WriteHeader(http.StatusBadRequest)
		// case errors.Contains(errorVal, errors.ErrUnsupportedContentType):
		// 	w.WriteHeader(http.StatusUnsupportedMediaType)

		// case errors.Contains(errorVal, errors.ErrMalformedEntity):
		// 	w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, errors.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		// case errors.Contains(errorVal, errors.ErrConflict):
		// 	w.WriteHeader(http.StatusConflict)

		// case errors.Contains(errorVal, errors.ErrScanMetadata),
		// 	errors.Contains(errorVal, errors.ErrSelectEntity):
		// 	w.WriteHeader(http.StatusUnprocessableEntity)

		case errors.Contains(errorVal, errors.ErrCreateEntity),
			errors.Contains(errorVal, errors.ErrUpdateEntity),
			errors.Contains(errorVal, errors.ErrViewEntity),
			errors.Contains(errorVal, errors.ErrRemoveEntity),
			errors.Contains(errorVal, errors.ErrConnect),
			errors.Contains(errorVal, errors.ErrDisconnect),
			errors.Contains(errorVal, errors.ErrMalformedEntity),
			errors.Contains(errorVal, errors.ErrAlreadyExists):
			//errors.Contains(errorVal, auth.ErrCreateGroup):
			w.WriteHeader(http.StatusBadRequest)

		case errors.Contains(errorVal, errors.ErrWrongPassword):
			w.WriteHeader(http.StatusBadRequest)

		case errors.Contains(errorVal, io.ErrUnexpectedEOF),
			errors.Contains(errorVal, io.EOF):
			w.WriteHeader(http.StatusBadRequest)

		case errors.Contains(errorVal, errors.ErrCreateUUID):
			w.WriteHeader(http.StatusInternalServerError)

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if errorVal.Msg() != "" {
			if err := json.NewEncoder(w).Encode(errorRes{Err: errorVal.Msg()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
