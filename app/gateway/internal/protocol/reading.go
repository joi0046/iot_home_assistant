package protocol

import "time"

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
	temperature *float64,
	humidity *float64,
	lux *float64,
) *Reading {
	data := make(map[string]Value)

	if temperature != nil {
		data["temperature"] = Value{
			Value: *temperature,
			Unit:  "°C",
		}
	}

	if humidity != nil {
		data["humidity"] = Value{
			Value: *humidity,
			Unit:  "%",
		}
	}

	if lux != nil {
		data["illuminance"] = Value{
			Value: *lux,
			Unit:  "lx",
		}
	}

	return &Reading{
		Version: "1",
		Device: Device{
			ID:   deviceID,
			Type: deviceType,
		},
		Data:      data,
		Timestamp: time.Now(),
	}
}
