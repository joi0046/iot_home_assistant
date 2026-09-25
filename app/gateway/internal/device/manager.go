package device

import (
	"gateway/internal/driver"
)

type Manager struct {
	sensors []driver.SensorDriver
}

func NewManager() *Manager {
	return &Manager{
		sensors: []driver.SensorDriver{},
	}
}

func (m *Manager) Register(sensor driver.SensorDriver) {
	m.sensors = append(m.sensors, sensor)
}

func (m *Manager) Sensors() []driver.SensorDriver {
	return m.sensors
}

func (m *Manager) DiscoverAndRegister(
	info driver.DeviceInfo,
	registry *driver.Registry,
) bool {
	sensor := registry.Find(info)

	if sensor == nil {
		return false
	}

	m.Register(sensor)
	return true
}
