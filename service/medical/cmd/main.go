package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"intelligent_guidance_system_v2/service/medical/internal/server"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf ./configs/config.yaml")
}

func main() {
	flag.Parse()

	logger := log.NewStdLogger(os.Stdout)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var cfg server.Config
	if err := c.Scan(&cfg); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(&cfg, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sig:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := app.Run(ctx); err != nil {
		log.Errorf("app run error: %v", err)
	}

	time.Sleep(1 * time.Second)
	log.Info("medical service exited")
}

func newApp(gs *grpc.Server, hs *http.Server, cfg *server.Config, logger log.Logger) *kratos.App {
	return kratos.New(
		kratos.ID("medical-service"),
		kratos.Name(cfg.Server.Name),
		kratos.Version("v1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}