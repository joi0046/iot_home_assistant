package mpu6050

import (
	"fmt"

	"gateway/internal/driver"
)

const (
	WHO_AM_I   = 0x75
	MPU6050_ID = 0x68
)

type I2C interface {
	SetAddress(address uint8) error
	ReadRegister(register uint8, length int) ([]byte, error)
}

type Sensor struct {
	bus     I2C
	address uint8
}

func New(bus I2C) *Sensor {
	return &Sensor{bus: bus}
}

func (s *Sensor) Name() string {
	return "MPU6050"
}

func (s *Sensor) Detect(info driver.DeviceInfo) bool {
	if info.Address != 0x68 && info.Address != 0x69 {
		return false
	}

	if err := s.bus.SetAddress(info.Address); err != nil {
		return false
	}

	data, err := s.bus.ReadRegister(WHO_AM_I, 1)
	if err != nil || len(data) != 1 {
		return false
	}

	if data[0] != MPU6050_ID {
		return false
	}

	s.address = info.Address
	return true
}

func (s *Sensor) Read() (driver.Reading, error) {
	if err := s.bus.SetAddress(s.address); err != nil {
		return nil, err
	}

	data, err := s.bus.ReadRegister(0x3B, 14)
	if err != nil {
		return nil, fmt.Errorf("read MPU6050 data: %w", err)
	}

	if len(data) != 14 {
		return nil, fmt.Errorf("unexpected data length: %d", len(data))
	}

	// 後でここに加速度・ジャイロ・温度の変換を入れる。
	return driver.Reading{}, nil
}
