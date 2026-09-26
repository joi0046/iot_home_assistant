package hdc1000

import (
	"testing"

	"gateway/internal/driver"
)

type fakeI2C struct {
	registers map[uint8][]byte
	readData  []byte
	address   uint8
}

func (f *fakeI2C) SetAddress(address uint8) error {
	f.address = address
	return nil
}

func (f *fakeI2C) Write(data []byte) error {
	return nil
}

func (f *fakeI2C) Read(length int) ([]byte, error) {
	return f.readData[:length], nil
}

func (f *fakeI2C) ReadRegister(
	register uint8,
	length int,
) ([]byte, error) {
	return f.registers[register][:length], nil
}

func TestDetect(t *testing.T) {
	bus := &fakeI2C{
		registers: map[uint8][]byte{
			ManufacturerRegister: {0x54, 0x49},
			DeviceRegister:       {0x10, 0x00},
		},
	}

	sensor := New(bus)

	info := driver.DeviceInfo{
		Bus:     "i2c-1",
		Address: 0x40,
	}

	if !sensor.Detect(info) {
		t.Error("Detect() = false, want true")
	}
}

func TestDetectUnsupportedAddress(t *testing.T) {
	bus := &fakeI2C{
		registers: map[uint8][]byte{
			ManufacturerRegister: {0x54, 0x49},
			DeviceRegister:       {0x10, 0x00},
		},
	}

	sensor := New(bus)

	info := driver.DeviceInfo{
		Bus:     "i2c-1",
		Address: 0x76,
	}

	if sensor.Detect(info) {
		t.Error("Detect() = true, want false")
	}
}

func TestRead(t *testing.T) {
	bus := &fakeI2C{
		readData: []byte{
			0x80, 0x00,
			0x80, 0x00,
		},
	}

	sensor := New(bus)

	sensor.address = 0x40

	reading, err := sensor.Read()
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if reading["temperature"] != 62.5 {
		t.Errorf(
			"temperature = %f, want 62.5",
			reading["temperature"],
		)
	}

	if reading["humidity"] != 50.0 {
		t.Errorf(
			"humidity = %f, want 50.0",
			reading["humidity"],
		)
	}
}
