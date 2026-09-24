package main

import (
	"fmt"
	"gateway/drivers/bh1750"
	"gateway/drivers/bme280"
	"gateway/internal/device"
	"gateway/internal/discovery"
)

func main() {
	discovery := discovery.New()
	discovery.Scan()

	manager := device.NewManager()
	manager.Register(&bme280.BME280{})
	manager.Register(&bh1750.BH1750{})

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
