package device

import (
	"gateway/internal/driver"
)

type Manager struct {
	devices []Device
}

func NewManager() *Manager {
	return &Manager{
		devices: []Device{},
	}
}

func (m *Manager) Register(device Device) {
	m.devices = append(m.devices, device)
}

func (m *Manager) Devices() []Device {
	return m.devices
}

func (m *Manager) DiscoverAndRegister(
	info driver.DeviceInfo,
	registry *driver.Registry,
) bool {
	sensor := registry.Find(info)

	if sensor == nil {
		return false
	}

	device := Device{
		Driver:  sensor,
		Bus:     "i2c-1",
		Address: info.Address,
	}

	m.Register(device)
	return true
}
