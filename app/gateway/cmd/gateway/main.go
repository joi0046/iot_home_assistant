package main

import (
	// packages
	"encoding/json"
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
	"gateway/internal/protocol"
)

func main() {
	// I²Cバスを開く
	bus, err := i2c.Open(1)
	if err != nil {
		fmt.Println("I²C error:", err)
		return
	}
	defer bus.Close()

	// I²C接続確認
	if err := bus.SetAddress(0x40); err != nil {
		fmt.Println("I²C address error:", err)
		return
	}

	fmt.Println("I²C connection OK")

	// HDC1000のID確認
	data, err := bus.ReadRegister(0xFE, 2)
	if err != nil {
		fmt.Println("FE read error:", err)
		return
	}
	fmt.Printf("FE: %02X %02X\n", data[0], data[1])

	data, err = bus.ReadRegister(0xFF, 2)
	if err != nil {
		fmt.Println("FF read error:", err)
		return
	}
	fmt.Printf("FF: %02X %02X\n", data[0], data[1])

	// MQTT接続
	mqttClient, err := mqtt.New("tcp://localhost:1883")
	if err != nil {
		fmt.Println("MQTT error:", err)
		return
	}
	defer mqttClient.Close()

	// センサードライバーを登録
	registry := driver.NewRegistry(
		&bme280.Sensor{},
		hdc1000.New(bus),
	)

	// Discovery
	discovery, err := discovery.New(registry)
	if err != nil {
		fmt.Println("Discovery error:", err)
		return
	}

	// Device Manager
	manager := device.NewManager()

	// デバイス検出
	devices := discovery.Scan()

	for _, info := range devices {
		if manager.DiscoverAndRegister(info, registry) {
			fmt.Printf(
				"Registered: %s/0x%02X\n",
				info.Bus,
				info.Address,
			)
		}
	}

	// 5秒ごとにセンサーを読み取る
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

			// 外部向けReadingへ変換
			reading := protocol.NewReading(
				device.Driver.Name(),
				fmt.Sprintf(
					"%s/0x%02X",
					device.Bus,
					device.Address,
				),
				value.Temperature,
				value.Humidity,
				value.Lux,
			)

			// JSONへ変換
			payload, err := json.Marshal(reading)
			if err != nil {
				fmt.Println("JSON error:", err)
				continue
			}

			// MQTTへ送信
			err = mqttClient.Publish(
				"gateway/readings",
				string(payload),
			)
			if err != nil {
				fmt.Println("MQTT publish error:", err)
				continue
			}

			fmt.Println("Published:", string(payload))

			// 人間向け表示
			if reading.Temperature != nil {
				fmt.Printf(
					"Temperature: %.2f °C\n",
					*reading.Temperature,
				)
			}

			if reading.Humidity != nil {
				fmt.Printf(
					"Humidity: %.2f %%\n",
					*reading.Humidity,
				)
			}

			if reading.Lux != nil {
				fmt.Printf(
					"Lux: %.2f lx\n",
					*reading.Lux,
				)
			}
		}
	}
}
