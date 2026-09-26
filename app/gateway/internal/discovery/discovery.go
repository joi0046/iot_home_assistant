package discovery

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"gateway/internal/driver"
)

type Scanner interface {
	Scan(busNumber string) ([]driver.DeviceInfo, error)
}

type I2CScanner struct{}

func (s *I2CScanner) Scan(busNumber string) ([]driver.DeviceInfo, error) {
	cmd := exec.Command("i2cdetect", "-y", busNumber)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("i2cdetect: %w", err)
	}

	devices := make([]driver.DeviceInfo, 0)

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) < 2 {
			continue
		}

		rowText := strings.TrimSuffix(fields[0], ":")

		row, err := strconv.ParseUint(rowText, 16, 8)
		if err != nil {
			continue
		}

		for column, value := range fields[1:] {
			if value == "--" || value == "UU" {
				continue
			}

			address, err := strconv.ParseUint(value, 16, 8)
			if err != nil {
				continue
			}

			expected := row + uint64(column)

			if address != expected {
				continue
			}

			devices = append(devices, driver.DeviceInfo{
				Bus:     "i2c-" + busNumber,
				Address: uint8(address),
			})
		}
	}

	return devices, nil
}

type Discovery struct {
	busNumber string
	scanner   Scanner
}

func New(scanner Scanner) (*Discovery, error) {
	return &Discovery{
		busNumber: "1",
		scanner:   scanner,
	}, nil
}

func (d *Discovery) Scan() []driver.DeviceInfo {
	fmt.Println("Scanning I²C bus...")

	devices, err := d.scanner.Scan(d.busNumber)
	if err != nil {
		fmt.Println("I²C scan error:", err)
		return nil
	}

	for _, device := range devices {
		fmt.Printf(
			"Found device: 0x%02X\n",
			device.Address,
		)
	}

	return devices
}
