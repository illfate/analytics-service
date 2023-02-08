package repository

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/illfate/analytics-service/internal/analytics"
)

type ClickHouse struct {
	conn driver.Conn
}

func NewClickHouse(conn driver.Conn) *ClickHouse {
	return &ClickHouse{conn: conn}
}

func (h *ClickHouse) Insert(ctx context.Context, events ...analytics.Event) error {
	batch, err := h.conn.PrepareBatch(ctx, `insert into events(
                   			client_time,
							device_id,
							device_os,
							session,
							sequence,
							event,
							param_int,
							param_str)
							`)
	if err != nil {
		return err
	}
	for _, event := range events {
		err := batch.Append(
			event.ClientTime,
			event.DeviceId,
			event.DeviceOs,
			event.Session,
			event.Sequence,
			event.Event,
			event.ParamInt,
			event.ParamStr,
		)
		if err != nil {
			return err
		}
	}
	return batch.Send()
}
