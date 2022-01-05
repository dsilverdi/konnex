package things

import (
	"context"
	"errors"
)

type Things struct {
	ID        string                 `json:"id"`
	ChannelID string                 `json:"channel_id"`
	Owner     string                 `json:"owner"`
	Name      string                 `json:"name"`
	MetaData  map[string]interface{} `json:"metadata"`
}

type ThingRepository interface {
	Insert(ctx context.Context, things Things) error
	GetAll(ctx context.Context) ([]Things, error)
}

//JSON Format Struct

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)
