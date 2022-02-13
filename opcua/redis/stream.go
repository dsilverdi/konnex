package redis

import (
	"context"
	"encoding/json"
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
	// thingUpdate     = thingPrefix + "update"
	// thingRemove     = thingPrefix + "remove"
	// thingConnect    = thingPrefix + "connect"
	// thingDisconnect = thingPrefix + "disconnect"

	// channelPrefix = "channel."
	// channelCreate = channelPrefix + "create"
	// channelUpdate = channelPrefix + "update"
	// channelRemove = channelPrefix + "remove"

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

				err = es.svc.CreateThing(ctx, cte.id, cte.serverUri, cte.nodeID)
			}
			if err != nil && err != errMetadataType {
				fmt.Println("Failed to handle event sourcing: ", err.Error())
				break
			}
			es.client.XAck(ctx, stream, group, msg.ID)
		}
	}
}

func decodeCreateThing(event map[string]interface{}) (createThingEvent, error) {
	var thingMetadata, channelMetadata map[string]interface{}

	thmeta := read(event, "thing_metadata", "{}")

	if err := json.Unmarshal([]byte(thmeta), &thingMetadata); err != nil {
		return createThingEvent{}, err
	}

	chmeta := read(event, "channel_metadata", "{}")

	if err := json.Unmarshal([]byte(chmeta), &channelMetadata); err != nil {
		return createThingEvent{}, err
	}

	cte := createThingEvent{
		id:   read(event, "id", ""),
		name: read(event, "name", ""),
	}

	NodeID, ok := thingMetadata[keyNodeID].(string)
	if !ok {
		return createThingEvent{}, errMetadataNodeID
	}

	ServerURI, ok := channelMetadata[keyServerURI].(string)
	if !ok {
		return createThingEvent{}, errMetadataServerURI
	}

	cte.nodeID = NodeID
	cte.serverUri = ServerURI
	return cte, nil
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}
	return val
}
