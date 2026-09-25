package i2c

import (
	"fmt"
	"os"
	"syscall"
)

const (
	I2CSlave = 0x0703
)

type Bus struct {
	file *os.File
}

func Open(busNumber int) (*Bus, error) {
	path := fmt.Sprintf("/dev/i2c-%d", busNumber)

	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open I2C bus: %w", err)
	}

	return &Bus{
		file: file,
	}, nil
}

func (b *Bus) Close() error {
	return b.file.Close()
}

func (b *Bus) SetAddress(address uint8) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		b.file.Fd(),
		I2CSlave,
		uintptr(address),
	)

	if errno != 0 {
		return fmt.Errorf("set I2C address 0x%02X: %w", address, errno)
	}

	return nil
}

func (b *Bus) ReadRegister(register uint8, length int) ([]byte, error) {
	_, err := b.file.Write([]byte{register})
	if err != nil {
		return nil, fmt.Errorf("write register: %w", err)
	}

	data := make([]byte, length)

	_, err = b.file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("read register: %w", err)
	}

	return data, nil
}
