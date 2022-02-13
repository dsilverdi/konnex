package gopcua

import (
	"context"
	"fmt"
	"konnex/opcua"
	"konnex/pkg/errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	opcuapkg "github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

// variable config

// const protocol = "opcua"
// const token = ""

var (
	// errNotFoundServerURI = errors.New("route map not found for Server URI")
	// errNotFoundNodeID    = errors.New("route map not found for Node ID")
	// errNotFoundConn      = errors.New("connection not found")

	errFailedConn = errors.New("failed to connect")
	// errFailedRead          = errors.New("failed to read")
	// errFailedParseInterval = errors.New("failed to parse subscription interval")
	// errFailedSub           = errors.New("failed to subscribe")
	// errFailedFindEndpoint  = errors.New("failed to find suitable endpoint")
	// errFailedFetchEndpoint = errors.New("failed to fetch OPC-UA server endpoints")
	// errFailedParseNodeID   = errors.New("failed to parse NodeID")
	// errFailedCreateReq     = errors.New("failed to create request")
	errFailedCreateMonitor = errors.New("failed to cretae Node Monitor")
	// errResponseStatus      = errors.New("response status not OK")
)

type Client struct {
	ctx      context.Context
	nodeData opcua.NodeDataRepository
}

func NewSubscriber(ctx context.Context, nodedataRepo opcua.NodeDataRepository) opcua.Subscriber {
	return &Client{
		ctx:      ctx,
		nodeData: nodedataRepo,
	}
}

func (cl *Client) Subscribe(_ context.Context, cfg opcua.Config, id string) error {
	opts := []opcuapkg.Option{
		opcuapkg.SecurityMode(ua.MessageSecurityModeNone),
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-signalCh
		println()
		cancel()
	}()

	oc := opcuapkg.NewClient(cfg.ServerURI, opts...)
	if err := oc.Connect(ctx); err != nil {
		return errors.Wrap(errFailedConn, err)
	}
	defer oc.Close()

	// i, err := strconv.Atoi(cfg.Interval)
	// if err != nil {
	// 	return errors.Wrap(errFailedParseInterval, err)
	// }

	mon, err := monitor.NewNodeMonitor(oc)
	if err != nil {
		return errors.Wrap(errFailedCreateMonitor, err)
	}

	mon.SetErrorHandler(func(_ *opcuapkg.Client, sub *monitor.Subscription, err error) {
		log.Printf("error: sub=%d err=%s", sub.SubscriptionID(), err.Error())
	})

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go cl.startCallbackSub(ctx, mon, time.Second, time.Millisecond*500, wg, id, cfg.NodeID)

	<-ctx.Done()
	wg.Wait()
	return nil
}

func (cl *Client) startCallbackSub(ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, wg *sync.WaitGroup, id string, nodes ...string) {
	fmt.Println("SUBS TO | ", nodes)
	sub, err := m.Subscribe(
		ctx,
		&opcuapkg.SubscriptionParameters{
			Interval: interval,
		},
		func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
			if msg.Error != nil {
				log.Printf("[callback] error=%s", msg.Error)
			} else {
				// log.Printf("[callback] node=%s value=%v", msg.NodeID, msg.Value.Value())
				cl.saveData(ctx, msg, id)
			}
			time.Sleep(lag)
		},
		nodes...)

	if err != nil {
		fmt.Print("error di sini startcallback")
		log.Fatal(err)
	}

	defer cl.cleanup(sub, wg)

	<-ctx.Done()
}

func (cl *Client) cleanup(sub *monitor.Subscription, wg *sync.WaitGroup) {
	log.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe(context.Background())
	wg.Done()
}

func (cl *Client) saveData(ctx context.Context, msg *monitor.DataChangeMessage, id string) {
	NewNodeData := &opcua.NodeData{
		Time:    time.Now(),
		ThingID: id,
	}

	switch msg.Value.Type() {
	case ua.TypeIDBoolean:
		NewNodeData.DataType = "boolean"
	case ua.TypeIDString, ua.TypeIDByteString:
		NewNodeData.DataType = "string"
	case ua.TypeIDDataValue:
		NewNodeData.DataType = "datavalue"
	case ua.TypeIDInt64, ua.TypeIDInt32, ua.TypeIDInt16:
		NewNodeData.DataType = "int"
	case ua.TypeIDUint64, ua.TypeIDUint32, ua.TypeIDUint16:
		NewNodeData.DataType = "uint"
	case ua.TypeIDFloat, ua.TypeIDDouble:
		NewNodeData.DataType = "float"
	case ua.TypeIDByte:
		NewNodeData.DataType = "byte"
	case ua.TypeIDDateTime:
		NewNodeData.DataType = "datetime"
	default:
		NewNodeData.DataType = "none"
	}

	msgVal := fmt.Sprintf("%v", msg.Value.Value())
	NewNodeData.Data = msgVal

	err := cl.nodeData.Save(ctx, NewNodeData)
	if err != nil {
		log.Printf("Error Save to DB | %s", err.Error())
	}
}
