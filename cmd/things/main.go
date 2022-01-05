package main

import (
	"flag"
	"fmt"
	"konnex/things"
	thingsapi "konnex/things/api"
	"konnex/things/sqldb"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"konnex/pkg/uuid"

	"github.com/go-kit/log"
	"github.com/jmoiron/sqlx"
)

const (
	DBHost     = "localhost"
	DBPort     = "3306"
	DBUser     = "root"
	DBPassword = "konnexthings"
	DBName     = "thingsdb"
	HTTPport   = ":8080"
)

type SysConfig struct {
	dbConfig sqldb.Config
	httpAddr string
}

func main() {
	cfg := loadConfig()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	db := connectDB(cfg)
	defer db.Close()

	svc := NewService(db)

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
		Host: DBHost,
		Port: DBPort,
		User: DBUser,
		Pass: DBPassword,
		Name: DBName,
	}

	var (
		httpAddr = flag.String("http.addr", HTTPport, "HTTP listen address")
	)
	flag.Parse()

	return SysConfig{
		dbConfig: myDB,
		httpAddr: *httpAddr,
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

func NewService(db *sqlx.DB) things.Service {
	database := sqldb.NewDatabase(db)

	ThingsRepository := sqldb.NewThingRepository(database)

	ChannelRepository := sqldb.NewChannelRepository(database)

	IDProvider := uuid.New()
	svc := things.New(ThingsRepository, ChannelRepository, IDProvider)
	return svc
}
