package bme280

import (
	"gateway/internal/driver"
)

type Sensor struct{}

func (s *Sensor) Name() string {
	return "BME280"
}

func (s *Sensor) Detect(info driver.DeviceInfo) bool {
	return info.Address == 0x76 || info.Address == 0x77
}

func (s *Sensor) Read() (driver.Reading, error) {
	return driver.Reading{}, nil
}
