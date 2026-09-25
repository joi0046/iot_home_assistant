package hdc1000

import (
	"fmt"

	"gateway/internal/driver"
	"gateway/internal/i2c"
)

const (
	ManufacturerRegister = 0xFE
	DeviceRegister       = 0xFF

	ManufacturerID uint16 = 0x5449
	DeviceID       uint16 = 0x1000
)

type Sensor struct {
	bus     *i2c.Bus
	address uint8
}

func New(bus *i2c.Bus) *Sensor {
	return &Sensor{
		bus: bus,
	}
}

func (s *Sensor) Name() string {
	return "HDC1000"
}

func (s *Sensor) Detect(info driver.DeviceInfo) bool {
	if info.Address < 0x40 || info.Address > 0x43 {
		return false
	}

	if err := s.bus.SetAddress(info.Address); err != nil {
		return false
	}

	manufacturer, err := s.bus.ReadRegister(ManufacturerRegister, 2)
	if err != nil {
		return false
	}

	device, err := s.bus.ReadRegister(DeviceRegister, 2)
	if err != nil {
		return false
	}

	manufacturerID := uint16(manufacturer[0])<<8 | uint16(manufacturer[1])
	deviceID := uint16(device[0])<<8 | uint16(device[1])

	return manufacturerID == ManufacturerID &&
		deviceID == DeviceID
}

func (s *Sensor) Read() (driver.Reading, error) {
	return driver.Reading{}, fmt.Errorf("HDC1000 Read is not implemented")
}
