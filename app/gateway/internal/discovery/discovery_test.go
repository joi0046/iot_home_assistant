package discovery

import (
	"errors"
	"testing"

	"gateway/internal/driver"
)

type testScanner struct {
	devices []driver.DeviceInfo
	err     error
}

func (s *testScanner) Scan(busNumber string) ([]driver.DeviceInfo, error) {
	return s.devices, s.err
}

func TestDiscoveryScan(t *testing.T) {
	scanner := &testScanner{
		devices: []driver.DeviceInfo{
			{
				Bus:     "i2c-1",
				Address: 0x40,
			},
			{
				Bus:     "i2c-1",
				Address: 0x76,
			},
		},
	}

	discovery, err := New(scanner)
	if err != nil {
		t.Fatal(err)
	}

	devices := discovery.Scan()

	if len(devices) != 2 {
		t.Fatalf(
			"len(devices) = %d, want 2",
			len(devices),
		)
	}

	if devices[0].Address != 0x40 {
		t.Errorf(
			"devices[0].Address = 0x%02X, want 0x40",
			devices[0].Address,
		)
	}

	if devices[1].Address != 0x76 {
		t.Errorf(
			"devices[1].Address = 0x%02X, want 0x76",
			devices[1].Address,
		)
	}
}

func TestDiscoveryScanError(t *testing.T) {
	scanner := &testScanner{
		err: errors.New("test error"),
	}

	discovery, err := New(scanner)
	if err != nil {
		t.Fatal(err)
	}

	devices := discovery.Scan()

	if devices != nil {
		t.Errorf(
			"Scan() = %v, want nil",
			devices,
		)
	}
}
