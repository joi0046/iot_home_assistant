package display

import (
	"fmt"

	"gateway/internal/driver"
)

func Print(value driver.Reading) {
	order := []string{
		"temperature",
		"humidity",
		"pressure",
		"illuminance",
		"co2",
		"voltage",
		"current",
		"power",
		"acceleration_x",
		"acceleration_y",
		"acceleration_z",
		"gyroscope_x",
		"gyroscope_y",
		"gyroscope_z",
	}

	for _, key := range order {
		if v, ok := value[key]; ok {
			switch key {
			case "temperature":
				fmt.Printf("Temperature: %.2f °C\n", v)

			case "humidity":
				fmt.Printf("Humidity: %.2f %%\n", v)

			case "pressure":
				fmt.Printf("Pressure: %.2f hPa\n", v)

			case "illuminance":
				fmt.Printf("Illuminance: %.2f lx\n", v)

			case "co2":
				fmt.Printf("CO2: %.0f ppm\n", v)

			case "voltage":
				fmt.Printf("Voltage: %.2f V\n", v)

			case "current":
				fmt.Printf("Current: %.2f A\n", v)

			case "power":
				fmt.Printf("Power: %.2f W\n", v)

			case "acceleration_x", "acceleration_y", "acceleration_z":
				fmt.Printf("%s: %.2f g\n", key, v)

			case "gyroscope_x", "gyroscope_y", "gyroscope_z":
				fmt.Printf("%s: %.2f °/s\n", key, v)
			}
		}
	}

}
