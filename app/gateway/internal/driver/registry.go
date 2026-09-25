package driver

type Registry struct {
	drivers []SensorDriver
}

func NewRegistry() *Registry {
	return &Registry{
		drivers: []SensorDriver{},
	}
}

func (r *Registry) Register(driver SensorDriver) {
	r.drivers = append(r.drivers, driver)
}

func (r *Registry) Find(info DeviceInfo) SensorDriver {
	for _, driver := range r.drivers {
		if driver.Detect(info) {
			return driver
		}
	}

	return nil
}
