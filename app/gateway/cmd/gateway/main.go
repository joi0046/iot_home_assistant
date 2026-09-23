package main

import (
	"fmt"
	"gateway/drivers/bh1750"
	"gateway/drivers/bme280"
	device "gateway/internal/devices"
	"gateway/internal/driver"
)

func main() {
	fmt.Println("IoT Gateway")
	fmt.Println("Gateway started")

	bme280dev := device.Device{
		ID:     "bme280_001",
		Name:   "BME280",
		Status: "online",
	}

	bh1750dev := device.Device{
		ID:     "bh1750_001",
		Name:   "BH1750",
		Status: "online",
	}

	fmt.Println("Device:", bme280dev)
	fmt.Println("Device:", bh1750dev)

	//sensor driverとして利用
	var sensor driver.SensorDriver = &bme280.BME280{}
	var bhSensor driver.SensorDriver = &bh1750.BH1750{}

	fmt.Println("Sensor:", sensor.Name())
	value, err := sensor.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Temperature:", value)

	fmt.Println("BH1750:", bhSensor.Name())
	bhValue, bhErr := bhSensor.Read()
	if bhErr != nil {
		fmt.Println("Error:", bhErr)
		return
	}

	fmt.Println("Light:", bhValue)
}
