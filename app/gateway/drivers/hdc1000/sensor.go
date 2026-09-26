package hdc1000

import (
	"time"

	"gateway/internal/driver"
	"gateway/internal/i2c"
)

type I2C interface {
	SetAddress(address uint8) error
	Write(data []byte) error
	Read(length int) ([]byte, error)
	ReadRegister(register uint8, length int) ([]byte, error)
}

const (
	ManufacturerRegister = 0xFE
	DeviceRegister       = 0xFF

	ManufacturerID uint16 = 0x5449
	DeviceID       uint16 = 0x1000
)

type Sensor struct {
	bus     I2C
	address uint8
}

func New(bus I2C) *Sensor {
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

	manufacturer, err := s.bus.ReadRegister(
		ManufacturerRegister,
		2,
	)
	if err != nil {
		return false
	}

	device, err := s.bus.ReadRegister(
		DeviceRegister,
		2,
	)
	if err != nil {
		return false
	}

	manufacturerID := uint16(manufacturer[0])<<8 |
		uint16(manufacturer[1])

	deviceID := uint16(device[0])<<8 |
		uint16(device[1])

	if manufacturerID == ManufacturerID &&
		deviceID == DeviceID {

		s.address = info.Address
		return true
	}

	return false
}

func (s *Sensor) Read() (driver.Reading, error) {
	if err := s.bus.SetAddress(s.address); err != nil {
		return driver.Reading{}, err
	}

	if err := s.bus.Write([]byte{0x00}); err != nil {
		return driver.Reading{}, err
	}

	time.Sleep(20 * time.Millisecond)

	data, err := s.bus.Read(4)
	if err != nil {
		return driver.Reading{}, err
	}

	rawTemperature := uint16(data[0])<<8 |
		uint16(data[1])

	rawHumidity := uint16(data[2])<<8 |
		uint16(data[3])

	temperature := float64(rawTemperature)/65536.0*165.0 - 40.0
	humidity := float64(rawHumidity) / 65536.0 * 100.0

	return driver.Reading{
		"temperature": temperature,
		"humidity":    humidity,
	}, nil
}

// 実際のI²C BusがI2Cインターフェースを満たすことを確認する。
// この宣言自体は実行時には何もしない。
var _ I2C = (*i2c.Bus)(nil)
