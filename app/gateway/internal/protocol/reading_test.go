package protocol

import (
	"testing"
	"time"

	"gateway/internal/driver"
)

func TestNewReading(t *testing.T) {
	data := driver.Reading{
		"temperature": 31.5,
		"humidity":    66.9,
	}

	reading := NewReading(
		"i2c-1/0x40",
		"HDC1000",
		data,
	)

	if reading.Version != "1" {
		t.Errorf(
			"Version = %q, want %q",
			reading.Version,
			"1",
		)
	}

	if reading.Device.ID != "i2c-1/0x40" {
		t.Errorf(
			"Device.ID = %q, want %q",
			reading.Device.ID,
			"i2c-1/0x40",
		)
	}

	if reading.Device.Type != "HDC1000" {
		t.Errorf(
			"Device.Type = %q, want %q",
			reading.Device.Type,
			"HDC1000",
		)
	}

	if reading.Data["temperature"].Value != 31.5 {
		t.Errorf(
			"temperature = %f, want %f",
			reading.Data["temperature"].Value,
			31.5,
		)
	}

	if reading.Data["temperature"].Unit != "°C" {
		t.Errorf(
			"temperature unit = %q, want %q",
			reading.Data["temperature"].Unit,
			"°C",
		)
	}

	if reading.Data["humidity"].Value != 66.9 {
		t.Errorf(
			"humidity = %f, want %f",
			reading.Data["humidity"].Value,
			66.9,
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
		t.Error("Timestamp is zero")
	}

	if time.Since(reading.Timestamp) > time.Second {
		t.Error("Timestamp is too old")
	}

}
