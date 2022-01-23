package main

import (
	"flag"
	"fmt"
	"konnex/things"
	thingsapi "konnex/things/api"
	rediscache "konnex/things/redis"
	"konnex/things/sqldb"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"konnex/pkg/uuid"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
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

	varHTTPport = ":8080"
)

type SysConfig struct {
	dbConfig  sqldb.Config
	redisURL  string
	redisPass string
	httpAddr  string
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

	svc := NewService(db, redisCl)

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

	redisURL := fmt.Sprintf("%s:%s", varRedisHost, varRedisPort)

	return SysConfig{
		dbConfig:  myDB,
		redisURL:  redisURL,
		redisPass: varRedisPass,
		httpAddr:  *httpAddr,
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

func NewService(db *sqlx.DB, redisCL *redis.Client) things.Service {
	database := sqldb.NewDatabase(db)

	ThingsRepository := sqldb.NewThingRepository(database)

	ChannelRepository := sqldb.NewChannelRepository(database)

	IDProvider := uuid.New()
	svc := things.New(ThingsRepository, ChannelRepository, IDProvider)

	svc = rediscache.NewEventStreamMiddleware(svc, redisCL)

	return svc
}
