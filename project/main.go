package main

import (
	//"master"
	//"./src/buttons"
	//"./src/definitions"
	//"./src/driver"
	"./src/storage"
	//"./src/master"
	//"./src/watchdog"

	//"./src/network"
	"fmt"
	//"os"
	//"time"
)

func main() {
	fmt.Printf("Test av lagring\n")
	storage.StoreInternalButtonPresses()
	storage.LoadInternalButtonPresses()
}
