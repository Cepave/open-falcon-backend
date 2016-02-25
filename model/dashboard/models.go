package dashboard

import (
	"time"
)

type Endpoint struct {
	ID       int64     `json:"id"`
	Endpoint string    `json:"endpoint"`
	Ts       int64     `json:"ts"`
	TCreate  time.Time `json:"-"`
	TModify  time.Time `json:"-"`
	Ipv4     string    `json:"-"`
}

type EndpointCounter struct {
	ID         int64     `json:"id"`
	EndpointID int64     `json:"endpoint_id"`
	Counter    string    `json:"counter"`
	Step       int64     `json:"step"`
	Type       string    `json:"type"`
	Ts         int64     `json:"ts"`
	TCreate    time.Time `json:"-"`
	TModify    time.Time `json:"-"`
}
