package driver

type DeviceInfo struct {
	Address uint8
}

type Sensor interface {
	Name() string
	Detect(info DeviceInfo) bool
	Read() (float64, error)
}
