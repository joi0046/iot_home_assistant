package mpu6050

import (
	"encoding/binary"
	"fmt"

	"gateway/internal/driver"
)

const (
	WHO_AM_I     = 0x75
	ACCEL_XOUT_H = 0x3B

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
	return &Sensor{
		bus: bus,
	}
}

func (s *Sensor) Initialize() error {
	if err := s.bus.SetAddress(s.address); err != nil {
		return err
	}

	// スリープ解除
	if err := s.bus.Write([]byte{0x6B, 0x00}); err != nil {
		return fmt.Errorf("wake MPU6050: %w", err)
	}

	return nil
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

	data, err := s.bus.ReadRegister(ACCEL_XOUT_H, 14)
	if err != nil {
		return nil, fmt.Errorf("read MPU6050 data: %w", err)
	}

	if len(data) != 14 {
		return nil, fmt.Errorf("unexpected data length: %d", len(data))
	}

	rawAccelX := int16(binary.BigEndian.Uint16(data[0:2]))
	rawAccelY := int16(binary.BigEndian.Uint16(data[2:4]))
	rawAccelZ := int16(binary.BigEndian.Uint16(data[4:6]))

	rawTemperature := int16(binary.BigEndian.Uint16(data[6:8]))

	rawGyroX := int16(binary.BigEndian.Uint16(data[8:10]))
	rawGyroY := int16(binary.BigEndian.Uint16(data[10:12]))
	rawGyroZ := int16(binary.BigEndian.Uint16(data[12:14]))

	// 加速度: ±2g → 16384 LSB/g
	accelerationX := float64(rawAccelX) / 16384.0
	accelerationY := float64(rawAccelY) / 16384.0
	accelerationZ := float64(rawAccelZ) / 16384.0

	// ジャイロ: ±250°/s → 131 LSB/(°/s)
	gyroscopeX := float64(rawGyroX) / 131.0
	gyroscopeY := float64(rawGyroY) / 131.0
	gyroscopeZ := float64(rawGyroZ) / 131.0

	// MPU6050内部温度
	temperature := float64(rawTemperature)/340.0 + 36.53

	return driver.Reading{
		"acceleration_x": accelerationX,
		"acceleration_y": accelerationY,
		"acceleration_z": accelerationZ,
		"gyroscope_x":    gyroscopeX,
		"gyroscope_y":    gyroscopeY,
		"gyroscope_z":    gyroscopeZ,
		"temperature":    temperature,
	}, nil
}
