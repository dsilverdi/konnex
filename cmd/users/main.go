package main

import (
	"flag"
	"fmt"
	"konnex/pkg/uuid"
	"konnex/users"
	userapi "konnex/users/api"
	"konnex/users/data"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/jmoiron/sqlx"
)

const (
	varDBHost     = "localhost"
	varDBPort     = "3306"
	varDBUser     = "root"
	varDBPassword = "konnexusers"
	varDBName     = "usersdb"

	// varRedisHost = "konnex-redis"
	// varRedisPort = "6379"
	// varRedisPass = ""

	varHTTPport = ":8081"
)

type SysConfig struct {
	dbConfig data.Config
	// redisURL  string
	// redisPass string
	httpAddr string
}

func main() {
	cfg := loadConfig()

	fmt.Println("PRINT THIS")

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
		h = userapi.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
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
	myDB := data.Config{
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

	// redisURL := fmt.Sprintf("%s:%s", varRedisHost, varRedisPort)

	return SysConfig{
		dbConfig: myDB,
		// redisURL:  redisURL,
		// redisPass: varRedisPass,
		httpAddr: *httpAddr,
	}
}

func connectDB(cfg SysConfig) *sqlx.DB {
	db, err := data.Connect(cfg.dbConfig)
	if err != nil {
		fmt.Println("error connecting db ", err)
		os.Exit(1)
	}

	return db
}

func NewService(db *sqlx.DB) users.Service {
	database := data.NewDatabase(db)

	UsersRepository := data.NewUsersRepository(database)

	AuthRepository := data.NewAuthRepository(database)

	IDProvider := uuid.New()
	svc := users.New(UsersRepository, AuthRepository, IDProvider)

	return svc
}
