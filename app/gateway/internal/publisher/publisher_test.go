package publisher

import (
	"encoding/json"
	"testing"
	"time"

	"gateway/internal/device"
	"gateway/internal/driver"
)

type fakeMQTT struct {
	topic   string
	payload string
}

func (f *fakeMQTT) Publish(topic string, payload string) error {
	f.topic = topic
	f.payload = payload
	return nil
}

func TestPublish(t *testing.T) {
	mqtt := &fakeMQTT{}
	publisher := New(mqtt)

	sensor := &testSensorDriver{
		name: "HDC1000",
	}

	device := device.Device{
		Driver:  sensor,
		Bus:     "i2c-1",
		Address: 0x40,
	}

	value := driver.Reading{
		"temperature": 25.5,
		"humidity":    60.0,
	}

	err := publisher.Publish(device, value)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	if mqtt.topic != ReadingTopic {
		t.Errorf(
			"topic = %q, want %q",
			mqtt.topic,
			ReadingTopic,
		)
	}

	var reading struct {
		Version string `json:"version"`

		Device struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"device"`

		Data map[string]struct {
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
		} `json:"data"`

		Timestamp time.Time `json:"timestamp"`
	}

	if err := json.Unmarshal(
		[]byte(mqtt.payload),
		&reading,
	); err != nil {
		t.Fatalf(
			"invalid JSON: %v\npayload: %s",
			err,
			mqtt.payload,
		)
	}

	if reading.Version != "1" {
		t.Errorf(
			"version = %q, want %q",
			reading.Version,
			"1",
		)
	}

	if reading.Device.ID != "i2c-1/0x40" {
		t.Errorf(
			"device.id = %q, want %q",
			reading.Device.ID,
			"i2c-1/0x40",
		)
	}

	if reading.Device.Type != "HDC1000" {
		t.Errorf(
			"device.type = %q, want %q",
			reading.Device.Type,
			"HDC1000",
		)
	}

	if reading.Data["temperature"].Value != 25.5 {
		t.Errorf(
			"temperature = %f, want 25.5",
			reading.Data["temperature"].Value,
		)
	}

	if reading.Data["temperature"].Unit != "°C" {
		t.Errorf(
			"temperature unit = %q, want %q",
			reading.Data["temperature"].Unit,
			"°C",
		)
	}

	if reading.Data["humidity"].Value != 60.0 {
		t.Errorf(
			"humidity = %f, want 60.0",
			reading.Data["humidity"].Value,
		)
	}

	if reading.Data["humidity"].Unit != "%" {
		t.Errorf(
			"humidity unit = %q, want %q",
			reading.Data["humidity"].Unit,
			"%",
		)
	}

	if reading.Timestamp.IsZero() {
		t.Error("timestamp is zero")
	}
}

type testSensorDriver struct {
	name string
}

func (d *testSensorDriver) Name() string {
	return d.name
}

func (d *testSensorDriver) Detect(info driver.DeviceInfo) bool {
	return false
}

func (d *testSensorDriver) Read() (driver.Reading, error) {
	return driver.Reading{}, nil
}
