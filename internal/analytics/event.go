package analytics

type Event struct {
	ClientTime string `json:"client_time"`
	DeviceId   string `json:"device_id"`
	DeviceOs   string `json:"device_os"`
	Session    string `json:"session"`
	Sequence   int    `json:"sequence"`
	Event      string `json:"event"`
	ParamInt   int    `json:"param_int"`
	ParamStr   string `json:"param_str"`
}
