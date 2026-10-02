package device

import (
	"testing"

	"gateway/internal/driver"
)

type testDriver struct {
	name    string
	address uint8
}

func (d *testDriver) Name() string {
	return d.name
}

func (d *testDriver) Detect(info driver.DeviceInfo) bool {
	return info.Address == d.address
}

func (d *testDriver) Read() (driver.Reading, error) {
	return driver.Reading{
		"temperature": 25.0,
	}, nil
}

func TestDiscoverAndRegister(t *testing.T) {
	testDriver := &testDriver{
		name:    "TestSensor",
		address: 0x40,
	}

	registry := driver.NewRegistry(testDriver)
	manager := NewManager()

	info := driver.DeviceInfo{
		Bus:     "i2c-1",
		Address: 0x40,
	}

	registered := manager.DiscoverAndRegister(info, registry)

	if !registered {
		t.Fatal("DiscoverAndRegister() = false, want true")
	}

	devices := manager.Devices()

	if len(devices) != 1 {
		t.Fatalf(
			"len(Devices()) = %d, want 1",
			len(devices),
		)
	}

	if devices[0].Driver.Name() != "TestSensor" {
		t.Errorf(
			"Driver.Name() = %q, want %q",
			devices[0].Driver.Name(),
			"TestSensor",
		)
	}

	if devices[0].Bus != "i2c-1" {
		t.Errorf(
			"Bus = %q, want %q",
			devices[0].Bus,
			"i2c-1",
		)
	}

	if devices[0].Address != 0x40 {
		t.Errorf(
			"Address = 0x%02X, want 0x40",
			devices[0].Address,
		)
	}

}

func TestDiscoverAndRegisterUnsupportedDevice(t *testing.T) {
	testDriver := &testDriver{
		name:    "TestSensor",
		address: 0x40,
	}

	registry := driver.NewRegistry(testDriver)
	manager := NewManager()

	info := driver.DeviceInfo{
		Bus:     "i2c-1",
		Address: 0x76,
	}

	registered := manager.DiscoverAndRegister(info, registry)

	if registered {
		t.Error("DiscoverAndRegister() = true, want false")
	}

	devices := manager.Devices()

	if len(devices) != 0 {
		t.Errorf(
			"len(Devices()) = %d, want 0",
			len(devices),
		)
	}

}
