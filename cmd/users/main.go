package main

import (
	"flag"
	"fmt"
	"konnex"
	"konnex/pkg/uuid"
	"konnex/users"
	usergrpcapi "konnex/users/api/grpc"
	userrestapi "konnex/users/api/rest"
	"konnex/users/data"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
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

	varHTTPport     = ":8081"
	varAuthGRPCport = ":9000"
	varAuthUrl      = "localhost:9000"
)

type SysConfig struct {
	dbConfig data.Config
	// redisURL  string
	// redisPass string
	httpAddr     string
	authGRPCport string
	authURL      string
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
		h = userrestapi.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
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

	go startGRPCServer(cfg.authGRPCport, svc, errs, logger)

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
		httpAddr:     *httpAddr,
		authGRPCport: varAuthGRPCport,
		authURL:      varAuthUrl,
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

func startGRPCServer(port string, svc users.Service, errs chan error, logger log.Logger) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Failed to listen on port %s: %s", port, err)
	}

	server := grpc.NewServer()

	konnex.RegisterAuthServiceServer(server, usergrpcapi.NewServer(svc))
	logger.Log("transport", "GRPC", "addr", port)
	errs <- server.Serve(listener)
}
