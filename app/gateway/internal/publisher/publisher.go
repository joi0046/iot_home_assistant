package publisher

import (
	"encoding/json"
	"fmt"

	"gateway/internal/device"
	"gateway/internal/driver"
	"gateway/internal/mqtt"
	"gateway/internal/protocol"
)

const ReadingTopic = "gateway/v1/readings"

type MQTT interface {
	Publish(topic string, payload string) error
}

type Publisher struct {
	mqtt MQTT
}

func New(mqttClient MQTT) *Publisher {
	return &Publisher{
		mqtt: mqttClient,
	}
}

func (p *Publisher) Publish(
	device device.Device,
	value driver.Reading,
) error {
	reading := protocol.NewReading(
		fmt.Sprintf(
			"%s/0x%02X",
			device.Bus,
			device.Address,
		),
		device.Driver.Name(),
		value,
	)

	payload, err := json.Marshal(reading)
	if err != nil {
		return fmt.Errorf("JSON error: %w", err)
	}

	if err := p.mqtt.Publish(
		ReadingTopic,
		string(payload),
	); err != nil {
		return fmt.Errorf("MQTT publish error: %w", err)
	}

	fmt.Println("Published:", string(payload))

	return nil
}

// 実際のMQTT ClientがMQTT interfaceを満たすことを確認する。
var _ MQTT = (*mqtt.Client)(nil)
