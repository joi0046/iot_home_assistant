package main

import (
	"fmt"
	"log"

	// drivers
	"gateway/drivers/bme280"

	// internal
	"gateway/internal/device"
	"gateway/internal/discovery"
	"gateway/internal/driver"
)

func main() {
	// Driver Registry
	registry := driver.NewRegistry(
		&bme280.Sensor{},
	)

	// Discovery
	d, err := discovery.New(registry)
	if err != nil {
		log.Fatal("Discovery error:", err)
	}

	devices := d.Scan()

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
