package things

import (
	"context"
	"errors"
)

type Things struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ThingRepository interface {
	Insert(ctx context.Context, things Things) error
}

//JSON Format Struct

type PostThingsRequest struct {
	Things Things
}

type PostThingsResponse struct {
	Err error `json:"err,omitempty"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)
