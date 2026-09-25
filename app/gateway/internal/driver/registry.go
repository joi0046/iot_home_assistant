package driver

type Registry struct {
	drivers []SensorDriver
}

func NewRegistry(drivers ...SensorDriver) *Registry {
	return &Registry{
		drivers: drivers,
	}
}

func (r *Registry) Find(info DeviceInfo) SensorDriver {
	for _, driver := range r.drivers {
		if driver.Detect(info) {
			return driver
		}
	}

	return nil
}
