package opcua

import (
	"context"
	"time"
)

type Node struct {
	ID        string
	ServerUri string
	NodeID    string
}

type NodeData struct {
	Time     time.Time
	ThingID  string
	Data     string
	DataType string
}

type NodeRepository interface {
	Save(context.Context, *Node) error
	ReadbyID(context.Context, string) (*Node, error)
	Delete(context.Context, string) error
}

type NodeDataRepository interface {
	Save(context.Context, *NodeData) error
	ReadbyID(context.Context, string) ([]NodeData, error)
}
