package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	I2CBus       int
	MQTTBroker   string
	ReadInterval time.Duration
}

func Load() Config {
	bus := 1
	if value := os.Getenv("IOT_I2C_BUS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			bus = parsed
		}
	}

	broker := os.Getenv("IOT_MQTT_BROKER")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	interval := 2 * time.Second
	if value := os.Getenv("IOT_READ_INTERVAL"); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			interval = parsed
		}
	}

	return Config{
		I2CBus:       bus,
		MQTTBroker:   broker,
		ReadInterval: interval,
	}
}
