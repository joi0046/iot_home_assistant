package driver

type SensorDriver interface {
	Name() string
	Read() (float64, error)
}
