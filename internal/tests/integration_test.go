package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/illfate/analytics-service/internal/analytics"
	"github.com/illfate/analytics-service/internal/handler"
	"github.com/illfate/analytics-service/internal/repository"
)

func TestCreateEvents(t *testing.T) {
	events := []analytics.Event{
		{
			ClientTime: time.Now(),
			DeviceId:   uuid.NewString(),
			DeviceOs:   "fsdfds",
			Session:    "fsdfsd",
			Sequence:   2,
			Event:      "fdsfs",
			ParamInt:   4,
			ParamStr:   "fsdfsd",
		},
		{
			ClientTime: time.Now(),
			DeviceId:   uuid.NewString(),
			DeviceOs:   "fsdfds",
			Session:    "fsdfsd",
			Sequence:   2,
			Event:      "fdsfs",
			ParamInt:   4,
			ParamStr:   "fsdfsd",
		},
	}
	var buf bytes.Buffer
	for _, event := range events {
		res, err := json.Marshal(event)
		require.NoError(t, err)
		buf.Write(res)
		buf.WriteString("\n")
	}
	makeTestReq(t, "/v1/events", http.MethodPost, &buf, http.StatusOK)

	assertEventsInClickHouse(t, events...)
}

func assertEventsInClickHouse(t *testing.T, events ...analytics.Event) {
	conn := mustConnectToClickHouse(t)
	for _, e := range events {
		shouldExistEventByDeviceID(t, conn, e.DeviceId)
	}
}

func makeTestReq(t *testing.T, urlPath string, method string, reader io.Reader, expStatus int) {
	clickHouse := mustConnectToClickHouse(t)

	house := repository.NewClickHouse(clickHouse)

	server := handler.NewServer(analytics.NewService(house), zaptest.NewLogger(t))
	testServer := httptest.NewServer(server)
	uri := testServer.URL + urlPath
	httpReq, err := http.NewRequest(method, uri, reader)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equalf(t, expStatus, resp.StatusCode, string(body))
}

func shouldExistEventByDeviceID(t *testing.T, conn driver.Conn, id string) {
	var count uint64
	err := conn.QueryRow(context.TODO(), `select count(*) from events where device_id = ?`, id).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), count)
}
