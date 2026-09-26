package main

import (
	// packages

	"fmt"
	"time"

	// drivers
	"gateway/drivers/bme280"
	"gateway/drivers/hdc1000"

	// internal
	"gateway/internal/device"
	"gateway/internal/discovery"
	"gateway/internal/driver"
	"gateway/internal/i2c"
	"gateway/internal/mqtt"
	"gateway/internal/publisher"
)

func main() {
	// I²Cバスを開く
	bus, err := i2c.Open(1)
	if err != nil {
		fmt.Println("I²C error:", err)
		return
	}
	defer bus.Close()

	// MQTT接続
	mqttClient, err := mqtt.New("tcp://localhost:1883")
	if err != nil {
		fmt.Println("MQTT error:", err)
		return
	}
	defer mqttClient.Close()

	// Publisher
	publisher := publisher.New(mqttClient)

	// センサードライバーを登録
	registry := driver.NewRegistry(
		&bme280.Sensor{},
		hdc1000.New(bus),
	)

	// Discovery
	scanner, err := discovery.New(registry)
	if err != nil {
		fmt.Println("Discovery error:", err)
		return
	}

	// Device Manager
	manager := device.NewManager()

	// デバイス検出
	devices := scanner.Scan()

	for _, info := range devices {
		if manager.DiscoverAndRegister(info, registry) {
			fmt.Printf(
				"Registered: %s/0x%02X\n",
				info.Bus,
				info.Address,
			)
		}
	}

	// 2秒ごとにセンサーを読み取る
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, device := range manager.Devices() {
			fmt.Println("Sensor:", device.Driver.Name())

			// センサー読み取り
			value, err := device.Driver.Read()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// MQTTへ送信
			if err := publisher.Publish(device, value); err != nil {
				fmt.Println("Publish error:", err)
				continue
			}

			// 人間向け表示
			for key, value := range value {
				switch key {
				case "temperature":
					fmt.Printf("Temperature: %.2f °C\n", value)

				case "humidity":
					fmt.Printf("Humidity: %.2f %%\n", value)

				case "illuminance":
					fmt.Printf("Illuminance: %.2f lx\n", value)

				case "pressure":
					fmt.Printf("Pressure: %.2f hPa\n", value)

				case "co2":
					fmt.Printf("CO2: %.0f ppm\n", value)

				default:
					fmt.Printf("%s: %.2f\n", key, value)
				}
			}
		}
	}
}
