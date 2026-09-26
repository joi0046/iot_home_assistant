package gateway

import (
	"fmt"
	"time"

	"gateway/drivers/bme280"
	"gateway/drivers/hdc1000"

	"gateway/internal/device"
	"gateway/internal/discovery"
	"gateway/internal/display"
	"gateway/internal/driver"
	"gateway/internal/i2c"
	"gateway/internal/mqtt"
	"gateway/internal/publisher"
)

type Gateway struct {
	bus       *i2c.Bus
	mqtt      *mqtt.Client
	publisher *publisher.Publisher
	manager   *device.Manager
}

func New() (*Gateway, error) {
	// I²Cバスを開く
	bus, err := i2c.Open(1)
	if err != nil {
		return nil, fmt.Errorf("I²C error: %w", err)
	}

	// MQTT接続
	mqttClient, err := mqtt.New("tcp://localhost:1883")
	if err != nil {
		bus.Close()
		return nil, fmt.Errorf("MQTT error: %w", err)
	}

	// Publisher
	publisher := publisher.New(mqttClient)

	// センサードライバーを登録
	registry := driver.NewRegistry(
		&bme280.Sensor{},
		hdc1000.New(bus),
	)

	// Discovery
	scanner, err := discovery.New(&discovery.I2CScanner{})
	if err != nil {
		mqttClient.Close()
		bus.Close()

		return nil, fmt.Errorf("Discovery error: %w", err)
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

	return &Gateway{
		bus:       bus,
		mqtt:      mqttClient,
		publisher: publisher,
		manager:   manager,
	}, nil
}

func (g *Gateway) Run() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, device := range g.manager.Devices() {
			fmt.Println("Sensor:", device.Driver.Name())

			// センサー読み取り
			value, err := device.Driver.Read()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// MQTTへ送信
			if err := g.publisher.Publish(device, value); err != nil {
				fmt.Println("Publish error:", err)
				continue
			}

			// 人間向け表示
			display.Print(value)
		}
	}
}

func (g *Gateway) Close() {
	g.mqtt.Close()
	g.bus.Close()
}
