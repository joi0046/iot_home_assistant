package discovery

import (
	"gateway/internal/driver"
)

type Discovery struct {
	addresses []uint8
}

func New() *Discovery {
	return &Discovery{
		addresses: []uint8{
			0x76,
		},
	}
}

func (d *Discovery) Scan() []driver.DeviceInfo {
	for _, address := range d.addresses {
		println("Found device:", address)
	}

	devices := make([]driver.DeviceInfo, 0)

	for _, address := range d.addresses {
		devices = append(devices, driver.DeviceInfo{
			Address: address,
		})
	}

	return devices
}
