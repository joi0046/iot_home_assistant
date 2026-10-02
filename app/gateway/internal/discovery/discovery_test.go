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

func TestParseOutput(t *testing.T) {
	output := `
     0 1 2 3 4 5 6 7 8 9 a b c d e f
00:                         -- -- -- -- -- -- -- --
10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
40: 40 -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
70: -- -- -- -- -- -- 76 -- -- -- -- -- -- -- -- --
`

	devices := parseOutput(output, "1")

	if len(devices) != 2 {
		t.Fatalf(
			"len(devices) = %d, want 2",
			len(devices),
		)
	}

	if devices[0].Bus != "i2c-1" {
		t.Errorf(
			"devices[0].Bus = %q, want %q",
			devices[0].Bus,
			"i2c-1",
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
