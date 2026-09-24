package discovery

import (
	"fmt"
)

type Discovery struct{}

func New() *Discovery {
	return &Discovery{}
}

func (d *Discovery) Scan() {
	fmt.Println("Scanning...")
	fmt.Println("No devices found.")
}
