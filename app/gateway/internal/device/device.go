package device

import "gateway/internal/driver"

type Device struct {
	Driver  driver.SensorDriver
	Bus     string
	Address uint8
}
