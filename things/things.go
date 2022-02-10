package things

import (
	"context"
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
	GetAll(ctx context.Context, owner string) ([]Things, error)
	GetSpecific(ctx context.Context, id, owner string) (*Things, error)
	Delete(ctx context.Context, id, owner string) error
}
