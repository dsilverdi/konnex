package opcua

import (
	"context"
	"fmt"
	"konnex/pkg/errors"
	"strings"
	"time"

	"strconv"
)

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type Service interface {
	// Browse OPCUA Node
	Browse(context.Context, string, string, string) ([]BrowsedNode, error)

	// Connect to OPCUA Server
	CreateThing(ctx context.Context, ThingsID, ServerURI, NodeID string) error

	// Subscribe From DB
	SubscribeWithDB(context.Context) error

	// Monitor OPC UA Data
	Monitor(context.Context, string) ([]MonitorData, error)
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
	Config      Config
	Browser     Browser
	Subscriber  Subscriber
	NodeRepo    NodeRepository
	MonitorRepo NodeDataRepository
}

type MonitorData struct {
	ID    string
	Value interface{}
	Time  time.Time
}

func NewService(cfg Config, browser Browser, sub Subscriber, noderepo NodeRepository, monitor NodeDataRepository) Service {
	return &opcuaService{
		Config:      cfg,
		Browser:     browser,
		Subscriber:  sub,
		NodeRepo:    noderepo,
		MonitorRepo: monitor,
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

func (svc opcuaService) SubscribeWithDB(ctx context.Context) error {
	nodes, err := svc.NodeRepo.ReadAll(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Saved Node Are | ", nodes)

	for _, node := range nodes {
		svc.Config.ServerURI = node.ServerUri
		svc.Config.NodeID = node.NodeID
		thingsid := node.ID
		go func() {
			if err := svc.Subscriber.Subscribe(ctx, svc.Config, thingsid); err != nil {
				fmt.Println("subscription failed", err)
			}
		}()
	}

	return nil
}

func (svc opcuaService) Monitor(ctx context.Context, id string) ([]MonitorData, error) {
	var datas []MonitorData

	fmt.Println("Service Monitor Jalan")

	results, err := svc.MonitorRepo.ReadbyID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	for i := range results {
		data := MonitorData{
			ID:   results[i].ThingID,
			Time: results[i].Time,
		}

		// Need to convert Data Value from string to desired Data Type Here
		switch results[i].DataType {
		case "boolean":
			data.Value = strings.ToLower(results[i].Data) == "true"
		case "string":
			data.Value = results[i].Data
		case "datavalue":
			data.Value = results[i].Data
		case "int":
			data.Value, _ = strconv.Atoi(results[i].Data)
		case "uint":
			data.Value, _ = strconv.Atoi(results[i].Data)
		case "float":
			data.Value, _ = strconv.ParseFloat(results[i].Data, 64)
		case "byte":
			data.Value = []byte(results[i].Data)
		case "datetime":
			data.Value = results[i].Data
		default:
			data.Value = results[i].Data
		}

		datas = append(datas, data)
	}

	return datas, nil
}
