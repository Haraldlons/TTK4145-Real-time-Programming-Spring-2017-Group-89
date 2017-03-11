package main

import (
	"./src/slave"
	// "./src/definitions"
	//"./src/driver"
	// "./src/elevator"
	// "./src/network"
	//"./src/buttons"
	"./src/driver"
	// "./src/storage"
	"./src/master"
	// "./src/watchdog"
	"fmt"
	// "log"
	// "os"
	// "os/exec"
	"time"
)

var delay = 50 * time.Millisecond
var elevatorActive = false

// var elevatorState = definitions.ElevatorState{2, 0}
var msg = make([]byte, 8)

func main() {
	fmt.Println("Main function started")
	// go slave.Run()
	// go master.Run()
	// go network.SetupNetwork()
	// Testchange

	go func() {
		stopSignal := 0
		for {
			stopSignal = driver.Elev_get_stop_signal()
			// fmt.Println("Stopsignal: ", stopSignal)
			if stopSignal != 0 {
				// setOrderOverNetwork(0)
				driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
				fmt.Println("Stopping program, with stop signal: ", stopSignal)
				fmt.Println("Another call to Elev_get_stop_signal(): ", driver.Elev_get_stop_signal())
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	slave.Run()
	master.Run()

	return

	// TestChange with haraldlons as user
	// Another testcommit with haraldlons@gmail.com as user.email

} //End mai
