package discovery

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"gateway/internal/driver"
)

type Discovery struct {
	busNumber string
}

func New() (*Discovery, error) {
	return &Discovery{
		busNumber: "1",
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

	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		for column, value := range fields[1:] {
			if value == "--" {
				continue
			}

			address, err := strconv.ParseUint(value, 16, 8)
			if err != nil {
				continue
			}

			// i2cdetectの行番号 + 列番号からアドレスを計算
			row, err := strconv.ParseUint(strings.TrimSuffix(fields[0], ":"), 16, 8)
			if err != nil {
				continue
			}

			address = row + uint64(column)

			fmt.Printf("Found device: 0x%02X\n", address)

			devices = append(devices, driver.DeviceInfo{
				Address: uint8(address),
			})
		}
	}

	return devices
}
