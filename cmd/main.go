package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/illfate/analytics-service/internal/analytics"
	"github.com/illfate/analytics-service/internal/handler"
	"github.com/illfate/analytics-service/internal/repository"
)

type Config struct {
	ClickHouseURL string `envconfig:"CLICK_HOUSE_URL" default:"127.0.0.1:9000"`
	ServerAddr    string `envconfig:"SERVICE_HOST" default:":11000"`
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return fmt.Errorf("failed to process env: %w", err)
	}
	ctx := context.Background()
	conn, err := connectToClickHouse(ctx, cfg.ClickHouseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to clickhouse: %w", err)
	}
	defer conn.Close()

	repo := repository.NewClickHouse(conn)
	service := analytics.NewService(repo)
	logger, err := zap.NewDevelopment() // TODO
	if err != nil {
		return fmt.Errorf("failed to setup zap: %w", err)
	}
	server := handler.NewServer(service, logger)

	err = runServer(service, cfg, server)
	if err != nil {
		return err
	}
	return nil
}

func runServer(service *analytics.Service, cfg Config, server *handler.Server) error {
	service.StartCreatingEvents(context.TODO())
	defer service.Close()

	httpServer := &http.Server{Addr: cfg.ServerAddr, Handler: server}
	exitCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var group errgroup.Group
	group.Go(func() error {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to listen and serve: %w", err)
		}
		return nil
	})
	group.Go(func() error {
		<-exitCtx.Done()
		return httpServer.Shutdown(context.TODO())
	})
	err := group.Wait()
	if err != nil {
		return err
	}
	return nil
}

func connectToClickHouse(ctx context.Context, addr string) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Debug: true,
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
