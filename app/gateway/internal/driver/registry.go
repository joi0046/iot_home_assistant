package driver

type Registry struct {
	sensors []Sensor
}

func NewRegistry(sensors ...Sensor) *Registry {
	return &Registry{
		sensors: sensors,
	}
}

func (r *Registry) Detect(info DeviceInfo) Sensor {
	for _, sensor := range r.sensors {
		if sensor.Detect(info) {
			return sensor
		}
	}

	return nil
}
