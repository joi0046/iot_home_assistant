package discovery

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"gateway/internal/driver"
)

type Discovery struct {
	busNumber string
	registry  *driver.Registry
}

func New(registry *driver.Registry) (*Discovery, error) {
	return &Discovery{
		busNumber: "1",
		registry:  registry,
	}, nil
}

func (d *Discovery) Scan() []driver.DeviceInfo {
	fmt.Println("Scanning I²C bus...")

	cmd := exec.Command("i2cdetect", "-y", d.busNumber)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("I²C scan error:", err)
		return nil
	}

	devices := make([]driver.DeviceInfo, 0)

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		// "00:", "10:", ... の行以外は無視
		if len(fields) < 2 {
			continue
		}

		rowText := strings.TrimSuffix(fields[0], ":")

		row, err := strconv.ParseUint(rowText, 16, 8)
		if err != nil {
			continue
		}

		for column, value := range fields[1:] {
			// デバイスが存在しない
			if value == "--" {
				continue
			}

			// UU = カーネルが使用中
			if value == "UU" {
				continue
			}

			address, err := strconv.ParseUint(value, 16, 8)
			if err != nil {
				continue
			}

			// 念のため i2cdetect の位置とも照合
			expected := row + uint64(column)

			if address != expected {
				continue
			}

			fmt.Printf("Found device: 0x%02X\n", address)

			devices = append(devices, driver.DeviceInfo{
				Address: uint8(address),
			})
		}
	}

	return devices
}
