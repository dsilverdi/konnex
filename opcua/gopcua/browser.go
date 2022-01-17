package gopcua

import (
	"context"
	"fmt"
	"konnex/opcua"
	"konnex/pkg/errors"
	"sync"

	opcuapkg "github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

const MaxChildren = 4 // max browsing node level

type Node struct {
	NodeID      *ua.NodeID
	NodeClass   ua.NodeClass
	BrowseName  string
	Description string
	AccessLevel ua.AccessLevelType
	Path        string
	DataType    string
	Writable    bool
	Unit        string
	Scale       string
	Min         string
	Max         string
}

type Browser struct {
	ctx context.Context
	// log logger.Logger
}

func NewBrowser(ctx context.Context) opcua.Browser {
	return Browser{
		ctx: ctx,
	}
}

func (b Browser) Browse(serveruri string, nodeid string) ([]opcua.BrowsedNode, error) {
	fmt.Println("browse jalan | ", serveruri)
	var nodelist []opcua.BrowsedNode

	opts := []opcuapkg.Option{
		opcuapkg.SecurityMode(ua.MessageSecurityModeNone),
	}

	uaClient := opcuapkg.NewClient(serveruri, opts...)
	if err := uaClient.Connect(b.ctx); err != nil {
		return nil, errors.Wrap(errFailedConn, err)
	}
	defer uaClient.Close()

	// nodes, err := browse(uaClient, nodeid, "", 0)
	// if err != nil {
	// 	return nil, err
	// }
	nodeCh := make(chan Node)
	var wg sync.WaitGroup

	wg.Add(1)
	go browse(uaClient, nodeid, "", 0, nodeCh, &wg)

	go func() {
		wg.Wait()
		close(nodeCh)
		fmt.Println("------- CLOSE NOW --------")
	}()

	var nodes []Node
	count := 0
	fmt.Println("--------- running read channel ---------")
	for ch := range nodeCh {
		count += 1
		nodes = append(nodes, ch)
	}

	fmt.Println("NODE SIZE | ", count)

	for _, n := range nodes {
		node := opcua.BrowsedNode{
			NodeID:      n.NodeID.String(),
			DataType:    n.DataType,
			Description: n.Description,
			Unit:        n.Unit,
			Scale:       n.Scale,
			BrowseName:  n.BrowseName,
		}
		nodelist = append(nodelist, node)
	}

	return nodelist, nil
}

func browse(cl *opcuapkg.Client, NodeID string, path string, level int, ch chan Node, wg *sync.WaitGroup) {
	defer wg.Done()

	if level > MaxChildren {
		// fmt.Println("max children")
		return
	}

	log := fmt.Sprintf("%s at %d", NodeID, level)
	fmt.Println("processing | ", log)

	nid, err := ua.ParseNodeID(NodeID)
	if err != nil {
		return
	}

	n := cl.Node(nid)

	attrs, err := n.Attributes(
		ua.AttributeIDNodeClass,
		ua.AttributeIDBrowseName,
		ua.AttributeIDDescription,
		ua.AttributeIDAccessLevel,
		ua.AttributeIDDataType,
	)
	if err != nil {
		return
	}

	nodeDef := new(Node)
	nodeDef.NodeID = nid

	switch err := attrs[0].Status; err {
	case ua.StatusOK:
		nodeDef.NodeClass = ua.NodeClass(attrs[0].Value.Int())
	default:
		return
	}

	switch err := attrs[1].Status; err {
	case ua.StatusOK:
		nodeDef.BrowseName = attrs[1].Value.String()
	default:
		return
	}

	switch err := attrs[2].Status; err {
	case ua.StatusOK:
		nodeDef.Description = attrs[2].Value.String()
	case ua.StatusBadAttributeIDInvalid:
		//pass
	default:
		return
	}

	switch err := attrs[3].Status; err {
	case ua.StatusOK:
		nodeDef.AccessLevel = ua.AccessLevelType(attrs[3].Value.Int())
		nodeDef.Writable = nodeDef.AccessLevel&ua.AccessLevelTypeCurrentWrite == ua.AccessLevelTypeCurrentWrite
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return
	}

	switch err := attrs[4].Status; err {
	case ua.StatusOK:
		switch v := attrs[4].Value.NodeID().IntID(); v {
		case id.DateTime:
			nodeDef.DataType = "time.Time"
		case id.Boolean:
			nodeDef.DataType = "bool"
		case id.SByte:
			nodeDef.DataType = "int8"
		case id.Int16:
			nodeDef.DataType = "int16"
		case id.Int32:
			nodeDef.DataType = "int32"
		case id.Byte:
			nodeDef.DataType = "byte"
		case id.UInt16:
			nodeDef.DataType = "uint16"
		case id.UInt32:
			nodeDef.DataType = "uint32"
		case id.UtcTime:
			nodeDef.DataType = "time.Time"
		case id.String:
			nodeDef.DataType = "string"
		case id.Float:
			nodeDef.DataType = "float32"
		case id.Double:
			nodeDef.DataType = "float64"
		default:
			nodeDef.DataType = attrs[4].Value.NodeID().String()
		}
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return
	}

	nodeDef.Path = join(path, nodeDef.BrowseName)

	if nodeDef.NodeClass == ua.NodeClassVariable {
		ch <- *nodeDef
		fmt.Println("----------- SEND TO CHANNEL ----------")
	}

	// ch <- *nodeDef

	browseChildren(cl, n, path, level, id.HasComponent, ch, wg)
	browseChildren(cl, n, path, level, id.Organizes, ch, wg)
	browseChildren(cl, n, path, level, id.HasProperty, ch, wg)

	// wg.Add(1)
	// go browsChildren(cl, n, path, level, id.HasComponent, ch, wg)

	// wg.Add(1)
	// go browsChildren(cl, n, path, level, id.Organizes, ch, wg)

	// wg.Add(1)
	// go browsChildren(cl, n, path, level, id.HasProperty, ch, wg)
}

func browseChildren(cl *opcuapkg.Client, node *opcuapkg.Node, path string, level int, refID uint32, ch chan Node, wg *sync.WaitGroup) {
	refs, err := node.ReferencedNodes(refID, ua.BrowseDirectionForward, ua.NodeClassAll, true)
	if err != nil {
		return
	}

	for _, r := range refs {
		wg.Add(1)
		go browse(cl, r.ID.String(), path, level+1, ch, wg)
	}
}

func join(a, b string) string {
	if a == "" {
		return b
	}
	return a + "." + b
}
