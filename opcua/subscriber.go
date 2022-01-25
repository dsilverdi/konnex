package opcua

import "context"

type Subscriber interface {
	Subscribe(context.Context, Config) error
}
