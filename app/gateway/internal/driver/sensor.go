package driver

type Reading struct {
	Temperature *float64
	Humidity    *float64
	Lux         *float64
}

type SensorDriver interface {
	Name() string
	Detect(info DeviceInfo) bool
	Read() (Reading, error)
}
