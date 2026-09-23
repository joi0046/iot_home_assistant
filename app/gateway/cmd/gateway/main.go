package main

import (
	"fmt"
	"gateway/drivers/bme280"
	device "gateway/internal/devices"
	"gateway/internal/driver"
)

func main() {
	fmt.Println("IoT Gateway")
	fmt.Println("Gateway started")

	d := device.Device{
		ID:     "bme280_001",
		Name:   "BME280",
		Status: "online",
	}

	fmt.Println("Device:", d)

	//sensor driverとして利用
	var sensor driver.SensorDriver = &bme280.Sensor{}

	fmt.Println("Sensor:", sensor.Name())
	value, err := sensor.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Temperature:", value)
}
