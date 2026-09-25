package main

import (
	"fmt"
	"log"

	// drivers
	"gateway/drivers/bme280"
	"gateway/drivers/hdc1000"

	// internal
	"gateway/internal/device"
	"gateway/internal/discovery"
	"gateway/internal/driver"
	"gateway/internal/i2c"
)

func main() {
	//test
	bus, err := i2c.Open(1)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	if err := bus.SetAddress(0x40); err != nil {
		log.Fatal(err)
	}

	fmt.Println("I²C connection OK")

	data, err := bus.ReadRegister(0xFE, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Manufacturer ID: 0x%02X%02X\n", data[0], data[1])
	// Driver Registry
	registry := driver.NewRegistry(
		&bme280.Sensor{},
		&hdc1000.Sensor{},
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
