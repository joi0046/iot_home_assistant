package discovery

import (
	"fmt"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"

	"gateway/internal/driver"
)

type Discovery struct {
	bus i2c.Bus
}

func New() (*Discovery, error) {
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	bus, err := i2creg.Open("1")
	if err != nil {
		return nil, err
	}

	return &Discovery{
		bus: bus,
	}, nil
}

func (d *Discovery) Scan() []driver.DeviceInfo {
	devices := make([]driver.DeviceInfo, 0)

	fmt.Println("Scanning I²C bus...")

	for address := 0x03; address <= 0x77; address++ {
		dev := &i2c.Dev{
			Bus:  d.bus,
			Addr: uint16(address),
		}

		err := dev.Tx(nil, nil)
		if err == nil {
			fmt.Printf("Found device: 0x%02X\n", address)

			devices = append(devices, driver.DeviceInfo{
				Address: uint8(address),
			})
		}
	}

	return devices
}
