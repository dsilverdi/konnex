package main

import (
	"context"
	"flag"
	"fmt"
	"konnex/opcua"
	opcuapi "konnex/opcua/api"
	"konnex/opcua/data"
	"konnex/opcua/gopcua"
	rediscl "konnex/opcua/redis"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	varOPCIntervalMs = "1000"
	varOPCPolicy     = ""
	varOPCMode       = ""
	varOPCCertFile   = ""
	varOPCKeyFile    = ""
	varHTTPport      = ":8082"
	varRedisHost     = "konnex-redis"
	varRedisPort     = "6379"
	varESConsumer    = "opcua"
)

type Config struct {
	uaConfig   opcua.Config
	httpAddr   string
	redisURL   string
	redisPass  string
	esConsumer string
}

func main() {
	cfg := LoadConfig()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	redisCl := connectRedis(cfg.redisURL, cfg.redisPass)

	ctx := context.Background()
	dbConn, err := data.ConnectDB(ctx)
	if err != nil {
		fmt.Println("Error Connecting to Database | ", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	sub := NewSubscriber(ctx, dbConn)
	svc := NewService(ctx, sub, cfg.uaConfig, dbConn)

	var h http.Handler
	{
		h = opcuapi.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go SubscribeFromSavedData()
	go SubscribeEventStream(svc, redisCl, cfg.esConsumer)
	go func() {
		logger.Log("transport", "HTTP", "addr", cfg.httpAddr)
		errs <- http.ListenAndServe(cfg.httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}

func LoadConfig() Config {
	opcuaCfg := opcua.Config{
		Interval: varOPCIntervalMs,
		Policy:   varOPCPolicy,
		Mode:     varOPCMode,
		CertFile: varOPCCertFile,
		KeyFile:  varOPCKeyFile,
	}

	var (
		httpAddr = flag.String("http.addr", varHTTPport, "HTTP listen address")
	)
	flag.Parse()

	redisURL := fmt.Sprintf("%s:%s", varRedisHost, varRedisPort)

	return Config{
		uaConfig:   opcuaCfg,
		httpAddr:   *httpAddr,
		redisURL:   redisURL,
		redisPass:  "",
		esConsumer: varESConsumer,
	}
}

func NewService(ctx context.Context, sub opcua.Subscriber, uaConfig opcua.Config, conn *pgxpool.Pool) opcua.Service {
	nodeBrowser := gopcua.NewBrowser(ctx)
	noderepo := data.NewNodeRepo(conn)
	svc := opcua.NewService(uaConfig, nodeBrowser, sub, noderepo)
	return svc
}

func NewSubscriber(ctx context.Context, conn *pgxpool.Pool) opcua.Subscriber {
	nodeData := data.NewNodeDataRepo(conn)
	return gopcua.NewSubscriber(ctx, nodeData)
}
func connectRedis(url string, pass string) *redis.Client {
	fmt.Println("Connect to Redis | ", url)
	return redis.NewClient(&redis.Options{
		Addr: url,
		// Password: pass,
	})
}

func SubscribeFromSavedData() {
	nodes, err := data.ReadAll()
	if err != nil {
		fmt.Println("Error Reading Data | ", err)
	}

	fmt.Println("Saved Node are | ", nodes)
}

func SubscribeEventStream(svc opcua.Service, client *redis.Client, consumer string) {
	eventStream := rediscl.NewEventStream(svc, client, consumer)
	if err := eventStream.Subscribe(context.Background(), "konnex.things"); err != nil {
		fmt.Println("Error Reading EventStream Redis")
	}
}
