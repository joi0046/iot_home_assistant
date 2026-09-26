package driver

import "testing"

type testSensorDriver struct {
	name    string
	address uint8
}

func (d *testSensorDriver) Name() string {
	return d.name
}

func (d *testSensorDriver) Detect(info DeviceInfo) bool {
	return info.Address == d.address
}

func (d *testSensorDriver) Read() (Reading, error) {
	return Reading{
		"temperature": 25.0,
	}, nil
}

func TestRegistryFind(t *testing.T) {
	bme280 := &testSensorDriver{
		name:    "BME280",
		address: 0x76,
	}

	hdc1000 := &testSensorDriver{
		name:    "HDC1000",
		address: 0x40,
	}

	registry := NewRegistry(
		bme280,
		hdc1000,
	)

	tests := []struct {
		name       string
		info       DeviceInfo
		wantDriver SensorDriver
	}{
		{
			name: "BME280",
			info: DeviceInfo{
				Bus:     "i2c-1",
				Address: 0x76,
			},
			wantDriver: bme280,
		},
		{
			name: "HDC1000",
			info: DeviceInfo{
				Bus:     "i2c-1",
				Address: 0x40,
			},
			wantDriver: hdc1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := registry.Find(tt.info)

			if got != tt.wantDriver {
				t.Errorf(
					"Find() = %v, want %v",
					got.Name(),
					tt.wantDriver.Name(),
				)
			}
		})
	}
}

func TestRegistryFindUnsupported(t *testing.T) {
	bme280 := &testSensorDriver{
		name:    "BME280",
		address: 0x76,
	}

	registry := NewRegistry(bme280)

	info := DeviceInfo{
		Bus:     "i2c-1",
		Address: 0x40,
	}

	got := registry.Find(info)

	if got != nil {
		t.Errorf(
			"Find() = %v, want nil",
			got.Name(),
		)
	}
}
