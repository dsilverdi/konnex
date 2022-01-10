package main

import (
	"context"
	"flag"
	"fmt"
	"konnex/opcua"
	opcuapi "konnex/opcua/api"
	"konnex/opcua/gopcua"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
)

const (
	OPCIntervalMs = "1000"
	OPCPolicy     = ""
	OPCMode       = ""
	OPCCertFile   = ""
	OPCKeyFile    = ""
	HTTPport      = ":8082"
)

type Config struct {
	uaConfig opcua.Config
	httpAddr string
}

func main() {
	cfg := LoadConfig()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := context.Background()
	svc := NewService(ctx, cfg.uaConfig)

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

	go func() {
		logger.Log("transport", "HTTP", "addr", cfg.httpAddr)
		errs <- http.ListenAndServe(cfg.httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}

func LoadConfig() Config {
	opcuaCfg := opcua.Config{
		Interval: OPCIntervalMs,
		Policy:   OPCPolicy,
		Mode:     OPCMode,
		CertFile: OPCCertFile,
		KeyFile:  OPCKeyFile,
	}

	var (
		httpAddr = flag.String("http.addr", HTTPport, "HTTP listen address")
	)
	flag.Parse()

	return Config{
		uaConfig: opcuaCfg,
		httpAddr: *httpAddr,
	}
}

func NewService(ctx context.Context, uaConfig opcua.Config) opcua.Service {
	nodeBrowser := gopcua.NewBrowser(ctx)
	svc := opcua.NewService(uaConfig, nodeBrowser)
	return svc
}
