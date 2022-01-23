package redis

import (
	"context"
	"fmt"
	"konnex/things"

	"github.com/go-redis/redis/v8"
)

const (
	streamID  = "konnex.things"
	streamLen = 1000
)

type EventStream struct {
	svc    things.Service
	client *redis.Client
}

func NewEventStreamMiddleware(svc things.Service, client *redis.Client) things.Service {
	return &EventStream{
		svc:    svc,
		client: client,
	}
}

func (es *EventStream) CreateThings(ctx context.Context, t things.Things) (*things.Things, error) {
	th, err := es.svc.CreateThings(ctx, t)
	if err != nil {
		return th, err
	}

	ch, err := es.svc.GetSpecificChannel(ctx, th.ChannelID)
	if err != nil {
		return th, err
	}

	event := createThingEvent{
		id:              th.ID,
		owner:           th.Owner,
		name:            th.Name,
		thingMetadata:   th.MetaData,
		channelMetadata: ch.Metadata,
	}

	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}

	err = es.client.XAdd(ctx, record).Err()
	if err != nil {
		fmt.Println("REDIS ERROR | ", err)
	}

	fmt.Println("publish to redis")

	return th, nil
}

func (es *EventStream) GetThings(ctx context.Context) ([]things.Things, error) {

	return es.svc.GetThings(ctx)
}

func (es *EventStream) GetSpecificThing(ctx context.Context, id string) (*things.Things, error) {
	return es.svc.GetSpecificThing(ctx, id)
}

func (es *EventStream) DeleteThing(ctx context.Context, id string) error {
	return es.svc.DeleteThing(ctx, id)
}

// Channel Services
func (es *EventStream) CreateChannel(ctx context.Context, ch things.Channel) (*things.Channel, error) {
	return es.svc.CreateChannel(ctx, ch)
}

func (es *EventStream) GetChannels(ctx context.Context) ([]things.Channel, error) {
	return es.svc.GetChannels(ctx)
}

func (es *EventStream) GetSpecificChannel(ctx context.Context, id string) (*things.Channel, error) {
	return es.svc.GetSpecificChannel(ctx, id)
}

func (es *EventStream) DeleteChannel(ctx context.Context, id string) error {
	return es.svc.DeleteChannel(ctx, id)
}
