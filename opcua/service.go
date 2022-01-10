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
	Config  Config
	Browser Browser
}

func NewService(cfg Config, browser Browser) Service {
	return &opcuaService{
		Config:  cfg,
		Browser: browser,
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
