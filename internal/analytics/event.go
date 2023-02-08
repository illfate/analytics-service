package analytics

import "time"

type Event struct {
	ClientTime time.Time `json:"client_time"`
	DeviceId   string    `json:"device_id"`
	DeviceOs   string    `json:"device_os"`
	Session    string    `json:"session"`
	Sequence   int64     `json:"sequence"`
	Event      string    `json:"event"`
	ParamInt   int64     `json:"param_int"`
	ParamStr   string    `json:"param_str"`
}
