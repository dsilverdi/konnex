package opcua

import "context"

type EventStream interface {
	Subscribe(context.Context, string) error
}
