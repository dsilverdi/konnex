package redis

import "encoding/json"

const (
	thingPrefix = "thing."
	thingCreate = thingPrefix + "create"
	thingDelete = thingPrefix + "delete"

	channelPrefix = "channel."
	// channelCreate = channelPrefix + "create"
	channelDelete = channelPrefix + "delete"
)

type createThingEvent struct {
	id              string
	channelID       string
	owner           string
	name            string
	thingMetadata   map[string]interface{}
	channelMetadata map[string]interface{}
}

func (cte createThingEvent) Encode() map[string]interface{} {
	val := map[string]interface{}{
		"id":         cte.id,
		"channel_id": cte.channelID,
		"owner":      cte.owner,
		"operation":  thingCreate,
	}

	if cte.name != "" {
		val["name"] = cte.name
	}

	if cte.thingMetadata != nil {
		metadata, err := json.Marshal(cte.thingMetadata)
		if err != nil {
			return val
		}

		val["thing_metadata"] = string(metadata)
	}

	if cte.channelMetadata != nil {
		metadata, err := json.Marshal(cte.channelMetadata)
		if err != nil {
			return val
		}

		val["channel_metadata"] = string(metadata)
	}

	return val
}

type deleteThingEvent struct {
	id string
}

func (e deleteThingEvent) Encode() map[string]interface{} {
	val := map[string]interface{}{
		"id":        e.id,
		"operation": thingDelete,
	}

	return val
}

type deleteChannelEvent struct {
	id string
}

func (e deleteChannelEvent) Encode() map[string]interface{} {
	val := map[string]interface{}{
		"id":        e.id,
		"operation": channelDelete,
	}

	return val
}
