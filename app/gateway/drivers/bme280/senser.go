package bme280

type Sensor struct{}

func (s *Sensor) Name() string {
	return "BME280"
}

func (s *Sensor) Read() (float64, error) {
	return 24.3, nil
}
