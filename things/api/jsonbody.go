package api

import "konnex/things"

// type PostThingsRequest struct {
// 	Things things.Things
// }

type createThingsResponse struct {
	Things  things.Things `json:"things,omitempty"`
	Err     error         `json:"err,omitempty"`
	Message string        `json:"message,omitempty"`
}

type createThingsReq struct {
	ID        string                 `json:"id,omitempty"`
	ChannelID string                 `json:"channel_id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	MetaData  map[string]interface{} `json:"metadata,omitempty"`
}

type getThingsRes struct {
	Message string          `json:"message,omitempty"`
	Things  []things.Things `json:"things,omitempty"`
}

type getThingsReq struct {
	channelID string
}

// Channel API Body

type createChannelResponse struct {
	Channel things.Channel `json:"channel,omitempty"`
	Err     error          `json:"err,omitempty"`
	Message string         `json:"message,omitempty"`
}

type createChannelReq struct {
	ID       string                 `json:"id,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Type     string                 `json:"type,omitempty"`
	MetaData map[string]interface{} `json:"metadata,omitempty"`
}

type getChannelReq struct {
	Type      string `json:"type,omitempty"`
	ChannelID string `json:"id,omitempty"`
}

type getChannelResponse struct {
	Message  string           `json:"message,omitempty"`
	Channels []things.Channel `json:"channels,omitempty"`
}
