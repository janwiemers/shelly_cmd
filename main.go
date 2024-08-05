package main

import (
	"log"
)

func main() {
	shelly := NewRpcApi("192.168.178.74")
	// err := shelly.PowerOnWithTimer(0, 10)

	// if err != nil {
	// 	panic(err)
	// }

	response, err := shelly.SwitchReset(0)

	if err != nil {
		panic(err)
	}

	log.Printf("%v", response)
}
