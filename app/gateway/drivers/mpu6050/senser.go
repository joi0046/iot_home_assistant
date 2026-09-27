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
	Write(data []byte) error
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

	if err := s.bus.Write([]byte{0x6B, 0x00}); err != nil {
		return nil, fmt.Errorf("wake MPU6050: %w", err)
	}

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

	readInt16 := func(high, low byte) int16 {
		return int16(uint16(high)<<8 | uint16(low))
	}

	rawAccelX := readInt16(data[0], data[1])
	rawAccelY := readInt16(data[2], data[3])
	rawAccelZ := readInt16(data[4], data[5])
	rawTemp := readInt16(data[6], data[7])
	rawGyroX := readInt16(data[8], data[9])
	rawGyroY := readInt16(data[10], data[11])
	rawGyroZ := readInt16(data[12], data[13])

	// 初期設定のフルスケール:
	// 加速度 ±2g → 16384 LSB/g
	// ジャイロ ±250°/s → 131 LSB/(°/s)
	accelX := float64(rawAccelX) / 16384.0
	accelY := float64(rawAccelY) / 16384.0
	accelZ := float64(rawAccelZ) / 16384.0

	gyroX := float64(rawGyroX) / 131.0
	gyroY := float64(rawGyroY) / 131.0
	gyroZ := float64(rawGyroZ) / 131.0

	temperature := float64(rawTemp)/340.0 + 36.53

	return driver.Reading{
		"acceleration_x": accelX,
		"acceleration_y": accelY,
		"acceleration_z": accelZ,
		"gyroscope_x":    gyroX,
		"gyroscope_y":    gyroY,
		"gyroscope_z":    gyroZ,
		"temperature":    temperature,
	}, nil
}
