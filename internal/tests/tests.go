package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/stretchr/testify/require"
)

func mustConnectToClickHouse(t *testing.T) driver.Conn {
	addr := os.Getenv("TEST_CLICK_HOUSE_URL")
	if addr == "" {
		addr = "127.0.0.1:9000"
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
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
	require.NoError(t, err)
	err = conn.Ping(context.TODO())
	require.NoError(t, err)
	return conn
}
