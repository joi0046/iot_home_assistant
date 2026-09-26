package protocol

import (
	"time"

	"gateway/internal/driver"
)

type Value struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Device struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Data map[string]Value

type Reading struct {
	Version   string    `json:"version"`
	Device    Device    `json:"device"`
	Data      Data      `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func NewReading(
	deviceID string,
	deviceType string,
	data driver.Reading,
) *Reading {
	values := make(Data)

	for key, value := range data {
		values[key] = Value{
			Value: value,
			Unit:  unitFor(key),
		}
	}

	return &Reading{
		Version: "1",
		Device: Device{
			ID:   deviceID,
			Type: deviceType,
		},
		Data:      values,
		Timestamp: time.Now(),
	}
}

func unitFor(key string) string {
	switch key {
	case "temperature":
		return "°C"
	case "humidity":
		return "%"
	case "pressure":
		return "hPa"
	case "illuminance":
		return "lx"
	case "co2":
		return "ppm"
	case "voltage":
		return "V"
	case "current":
		return "A"
	case "power":
		return "W"
	default:
		return ""
	}
}
