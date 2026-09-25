package protocol

import "time"

type Reading struct {
	Device      string    `json:"device"`
	Address     string    `json:"address"`
	Temperature *float64  `json:"temperature,omitempty"`
	Humidity    *float64  `json:"humidity,omitempty"`
	Lux         *float64  `json:"lux,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewReading(
	device string,
	address string,
	temperature *float64,
	humidity *float64,
	lux *float64,
) Reading {
	return Reading{
		Device:      device,
		Address:     address,
		Temperature: temperature,
		Humidity:    humidity,
		Lux:         lux,
		Timestamp:   time.Now(),
	}
}
