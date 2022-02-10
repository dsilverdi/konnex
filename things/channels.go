package things

import "context"

type Channel struct {
	ID       string                 `json:"id"`
	Owner    string                 `json:"owner"`
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type ChannelRepository interface {
	Insert(ctx context.Context, ch Channel) error
	GetAll(ctx context.Context, owner string) ([]Channel, error)
	GetSpecific(ctx context.Context, id, owner string) (*Channel, error)
	Delete(ctx context.Context, id, owner string) error
}
