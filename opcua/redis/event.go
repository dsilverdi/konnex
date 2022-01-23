package redis

type createThingEvent struct {
	id        string
	name      string
	serverUri string
	nodeID    string
}
