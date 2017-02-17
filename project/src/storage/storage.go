package storage

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

//Internal button presses:
func StoreInternalButtonPresses() bool {

	return true
}

func LoadInternalButtonPresses() bool {
	return true
}

//External button presses:
func StoreExternalButtonPresses() bool {
	return true
}

func LoadExternalButtonPresses() bool {
	return true
}

func StoreOrders(elevatorNum int) bool {
	return true
}

func LoadOrders(elvatorNum int) bool {
	return true
}
