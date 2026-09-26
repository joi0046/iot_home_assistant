package main

import (
	"fmt"

	// gateway
	"gateway/internal/gateway"
)

func main() {
	gateway, err := gateway.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gateway.Close()

	gateway.Run()
}
