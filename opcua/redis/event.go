package redis

import "encoding/json"

type createThingEvent struct {
	id        string
	channelID string
	name      string
	serverUri string
	nodeID    string
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
		id:        read(event, "id", ""),
		channelID: read(event, "channel_id", ""),
		name:      read(event, "name", ""),
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

type deleteThingEvent struct {
	id string
}

func decodeDeleteThing(event map[string]interface{}) (deleteThingEvent, error) {
	dte := deleteThingEvent{
		id: read(event, "id", ""),
	}

	return dte, nil
}

type deleteChannelEvent struct {
	id string
}

func decodeDeleteChannel(event map[string]interface{}) (deleteChannelEvent, error) {
	dce := deleteChannelEvent{
		id: read(event, "id", ""),
	}

	return dce, nil
}
