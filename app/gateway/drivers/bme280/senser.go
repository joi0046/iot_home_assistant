package bme280

type BME280 struct{}

func (s *BME280) Name() string {
	return "BME280"
}

func (s *BME280) Read() (float64, error) {
	return 24.3, nil
}
