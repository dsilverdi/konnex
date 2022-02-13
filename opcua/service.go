package opcua

import (
	"context"
	"fmt"
	"konnex/pkg/errors"
)

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type Service interface {
	//Browse OPCUA Node
	Browse(context.Context, string, string, string) ([]BrowsedNode, error)

	//Connect to OPCUA Server
	CreateThing(ctx context.Context, ThingsID, ServerURI, NodeID string) error
}

type Config struct {
	ServerURI string
	NodeID    string
	Interval  string
	Policy    string
	Mode      string
	CertFile  string
	KeyFile   string
}

type opcuaService struct {
	Config     Config
	Browser    Browser
	Subscriber Subscriber
	NodeRepo   NodeRepository
}

func NewService(cfg Config, browser Browser, sub Subscriber, noderepo NodeRepository) Service {
	return &opcuaService{
		Config:     cfg,
		Browser:    browser,
		Subscriber: sub,
		NodeRepo:   noderepo,
	}
}

func (svc opcuaService) Browse(ctx context.Context, serveruri, namespace, identifier string) ([]BrowsedNode, error) {
	nodeID := fmt.Sprintf("%s;%s", namespace, identifier)

	nodes, err := svc.Browser.Browse(serveruri, nodeID)
	if err != nil {
		return nil, errors.Wrap(ErrNotFound, err)
	}

	return nodes, nil
}

func (svc opcuaService) CreateThing(ctx context.Context, ThingsID, ServerURI, NodeID string) error {
	fmt.Println("Got IoT Data Called From Redis | ", []string{ServerURI, NodeID})

	svc.Config.ServerURI = ServerURI
	svc.Config.NodeID = NodeID

	go func() {
		if err := svc.Subscriber.Subscribe(ctx, svc.Config, ThingsID); err != nil {
			fmt.Println("subscription failed", err)
		}
	}()

	NewNode := &Node{
		ID:        ThingsID,
		ServerUri: ServerURI,
		NodeID:    NodeID,
	}

	return svc.NodeRepo.Save(ctx, NewNode)
}
