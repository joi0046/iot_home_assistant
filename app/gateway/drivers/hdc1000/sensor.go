package hdc1000

import (
	"fmt"
	"gateway/internal/driver"
)

const (
	ManufacturerID uint16 = 0x5449
	DeviceID       uint16 = 0x1000

	ManufacturerRegister uint8 = 0xFE
	DeviceRegister       uint8 = 0xFF
)

type Sensor struct {
	Address uint8
}

func (s *Sensor) Name() string {
	return "HDC1000"
}

func (s *Sensor) Detect(info driver.DeviceInfo) bool {
	return info.Address == 0x40 ||
		info.Address == 0x41 ||
		info.Address == 0x42 ||
		info.Address == 0x43
}

func (s *Sensor) Read() (float64, error) {
	return 0, fmt.Errorf("HDC1000 Read is not implemented")
}
