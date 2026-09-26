package driver

type Reading map[string]float64

type SensorDriver interface {
	Name() string
	Detect(info DeviceInfo) bool
	Read() (Reading, error)
}
