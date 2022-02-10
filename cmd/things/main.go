package main

import (
	"flag"
	"fmt"
	"konnex"
	"konnex/things"
	thingsapi "konnex/things/api"
	rediscache "konnex/things/redis"
	"konnex/things/sqldb"
	authapi "konnex/users/api/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"konnex/pkg/uuid"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

const (
	varDBHost     = "localhost"
	varDBPort     = "3306"
	varDBUser     = "root"
	varDBPassword = "konnexthings"
	varDBName     = "thingsdb"

	varRedisHost = "konnex-redis"
	varRedisPort = "6379"
	varRedisPass = ""

	varHTTPport     = ":8080"
	varAuthGRPCport = ":9000"
	varAuthUrl      = "konnex-users:9000"
)

type SysConfig struct {
	dbConfig     sqldb.Config
	redisURL     string
	redisPass    string
	httpAddr     string
	authGRPCport string
	authURL      string
	authTimeout  time.Duration
}

func main() {
	cfg := loadConfig()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	redisCl := connectRedis(cfg.redisURL, cfg.redisPass)

	db := connectDB(cfg)
	defer db.Close()

	auth, close := InitAuthClient(cfg)
	if close != nil {
		defer close()
	}

	svc := NewService(db, redisCl, auth)

	var h http.Handler
	{
		h = thingsapi.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", cfg.httpAddr)
		errs <- http.ListenAndServe(cfg.httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}

func loadConfig() SysConfig {
	myDB := sqldb.Config{
		Host: varDBHost,
		Port: varDBPort,
		User: varDBUser,
		Pass: varDBPassword,
		Name: varDBName,
	}

	var (
		httpAddr = flag.String("http.addr", varHTTPport, "HTTP listen address")
	)
	flag.Parse()

	authTimeout, err := time.ParseDuration("1s")
	if err != nil {
		fmt.Printf("Invalid %s value: %s", "1s", err.Error())
		os.Exit(1)
	}

	redisURL := fmt.Sprintf("%s:%s", varRedisHost, varRedisPort)

	return SysConfig{
		dbConfig:     myDB,
		redisURL:     redisURL,
		redisPass:    varRedisPass,
		httpAddr:     *httpAddr,
		authGRPCport: varAuthGRPCport,
		authURL:      varAuthUrl,
		authTimeout:  authTimeout,
	}
}

func connectDB(cfg SysConfig) *sqlx.DB {
	db, err := sqldb.Connect(cfg.dbConfig)
	if err != nil {
		fmt.Println("error connecting db ", err)
		os.Exit(1)
	}

	return db
}

func connectRedis(url string, pass string) *redis.Client {
	fmt.Println("Connect to Redis | ", url)
	return redis.NewClient(&redis.Options{
		Addr: url,
		// Password: pass,
	})
}

func NewService(db *sqlx.DB, redisCL *redis.Client, auth konnex.AuthServiceClient) things.Service {
	database := sqldb.NewDatabase(db)

	ThingsRepository := sqldb.NewThingRepository(database)

	ChannelRepository := sqldb.NewChannelRepository(database)

	IDProvider := uuid.New()
	svc := things.New(ThingsRepository, ChannelRepository, IDProvider, auth)

	svc = rediscache.NewEventStreamMiddleware(svc, redisCL)

	return svc
}

func InitAuthClient(cfg SysConfig) (konnex.AuthServiceClient, func() error) {
	conn := connectAuthGRPC(cfg)
	return authapi.NewClient(conn, cfg.authTimeout), conn.Close
}

func connectAuthGRPC(cfg SysConfig) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(cfg.authURL, opts...)
	if err != nil {
		fmt.Printf("Failed to connect to auth service: %s", err)
		os.Exit(1)
	}

	return conn
}
