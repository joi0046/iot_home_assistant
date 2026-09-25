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
	"gateway/internal/protocol"
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

	fmt.Printf("FE: %02X %02X\n", data[0], data[1])

	data, err = bus.ReadRegister(0xFF, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("FF: %02X %02X\n", data[0], data[1])
	// Driver Registry
	registry := driver.NewRegistry(
		&bme280.Sensor{},
		hdc1000.New(bus),
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
	for _, device := range manager.Devices() {
		fmt.Println("Sensor:", device.Driver.Name())

		value, err := device.Driver.Read()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		reading := protocol.NewReading(
			device.Driver.Name(),
			fmt.Sprintf("%s/0x%02X", device.Bus, device.Address),
			value.Temperature,
			value.Humidity,
			value.Lux,
		)

		fmt.Printf("Reading: %+v\n", reading)

		if reading.Temperature != nil {
			fmt.Printf("Temperature: %.2f °C\n", *reading.Temperature)
		}

		if reading.Humidity != nil {
			fmt.Printf("Humidity: %.2f %%\n", *reading.Humidity)
		}

		if reading.Lux != nil {
			fmt.Printf("Lux: %.2f lx\n", *reading.Lux)
		}
	}
}
