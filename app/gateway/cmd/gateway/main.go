package main

import (
	"fmt"
	//drivers
	"gateway/drivers/bme280"
	//internal
	"gateway/internal/device"
	"gateway/internal/discovery"
	"gateway/internal/driver"
)

func main() {
	// Discovery
	d, err := discovery.New()
	if err != nil {
		fmt.Println("Discovery error:", err)
		return
	}

	devices := d.Scan()

	for _, info := range devices {
		fmt.Printf("Found: 0x%02X\n", info.Address)
	}

	// Driver Registry
	registry := driver.NewRegistry()
	registry.Register(&bme280.Sensor{})

	// Device Manager
	manager := device.NewManager()

	for _, info := range devices {
		if manager.DiscoverAndRegister(info, registry) {
			fmt.Printf("Registered: 0x%02X\n", info.Address)
		} else {
			fmt.Printf("Unknown device: 0x%02X\n", info.Address)
		}
	}

	// Sensor Read
	for _, sensor := range manager.Sensors() {
		fmt.Println("Sensor:", sensor.Name())

		value, err := sensor.Read()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Value:", value)
	}
}
