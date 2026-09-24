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
