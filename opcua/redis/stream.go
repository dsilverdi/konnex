package redis

import (
	"context"
	"errors"
	"fmt"
	"konnex/opcua"

	"github.com/go-redis/redis/v8"
)

const (
	keyType      = "opcua"
	keyNodeID    = "node"
	keyServerURI = "server-uri"

	group  = "konnex.opcua"
	stream = "konnex.things"

	thingPrefix = "thing."
	thingCreate = thingPrefix + "create"
	thingDelete = thingPrefix + "delete"

	channelPrefix = "channel."
	// channelCreate = channelPrefix + "create"
	channelDelete = channelPrefix + "delete"

	exists = "BUSYGROUP Consumer Group name already exists"
)

var (
	errMetadataType = errors.New("metadatada is not of type opcua")

	// errMetadataFormat = errors.New("malformed metadata")

	errMetadataServerURI = errors.New("ServerURI not found in channel metadatada")

	errMetadataNodeID = errors.New("NodeID not found in thing metadatada")
)

type eventStream struct {
	svc      opcua.Service
	client   *redis.Client
	consumer string
}

func NewEventStream(svc opcua.Service, cl *redis.Client, consumer string) opcua.EventStream {
	return &eventStream{
		svc:      svc,
		client:   cl,
		consumer: consumer,
	}
}

func (es *eventStream) Subscribe(ctx context.Context, stream string) error {
	err := es.client.XGroupCreateMkStream(ctx, stream, group, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	fmt.Println("Subscribe to EventStream | ", stream)

	for {
		streams, err := es.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: es.consumer,
			Streams:  []string{stream, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values

			var err error
			switch event["operation"] {
			case thingCreate:
				cte, e := decodeCreateThing(event)
				if e != nil {
					err = e
					break
				}

				err = es.svc.CreateThing(ctx, cte.id, cte.channelID, cte.serverUri, cte.nodeID)

			case thingDelete:
				dte, e := decodeDeleteThing(event)
				if e != nil {
					err = e
					break
				}

				err = es.svc.DeleteThing(ctx, dte.id)

			case channelDelete:
				cde, e := decodeDeleteChannel(event)
				if e != nil {
					err = e
					break
				}

				err = es.svc.DeleteChannel(ctx, cde.id)
			}
			if err != nil && err != errMetadataType {
				fmt.Println("Failed to handle event sourcing: ", err.Error())
				break
			}
			es.client.XAck(ctx, stream, group, msg.ID)
		}
	}
}
