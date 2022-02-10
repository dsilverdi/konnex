package api

// type PostThingsRequest struct {
// 	Things things.Things
// }

//
type HTTPResponse struct {
	Code    int         `json:"code,omitempty"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type getSpecificReq struct {
	Token string
	ID    string
}

//
type createThingsReq struct {
	Token     string
	ID        string                 `json:"id,omitempty"`
	ChannelID string                 `json:"channel_id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	MetaData  map[string]interface{} `json:"metadata,omitempty"`
}

type getThingsReq struct {
	Token     string
	channelID string
}

// Channel API Body
type createChannelReq struct {
	Token    string
	ID       string                 `json:"id,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Type     string                 `json:"type,omitempty"`
	MetaData map[string]interface{} `json:"metadata,omitempty"`
}

type getChannelReq struct {
	Token     string
	Type      string `json:"type,omitempty"`
	ChannelID string `json:"id,omitempty"`
}
