package main

import (
	"./src/definitions"
	"./src/slave"
	//"./src/driver"
	// "./src/elevator"
	"./src/network"
	//"./src/buttons"
	"./src/driver"
	// "./src/storage"
	"./src/master"
	// "./src/watchdog"
	// "fmt"
	// "log"
	// "os"
	// "os/exec"
	// "net"
	"time"
)

func main() {

	go func() {
		stopSignal := 0
		// time.Sleep(500 * time.Millisecond)
		for {
			stopSignal = driver.Elev_get_stop_signal()
			// fmt.Println("Stopsignal: ", stopSignal)
			if stopSignal != 0 {
				// setOrderOverNetwork(0)
				driver.Elev_set_motor_direction(definitions.DIR_STOP)
				// fmt.Println("Stopping program, with stop signal: ", stopSignal)
				// fmt.Println("Another call to Elev_get_stop_signal(): ", driver.Elev_get_stop_signal())
				time.Sleep(100 * time.Millisecond)
				// return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if network.CheckIfMasterAlreadyExist() {
				slave.Run()
				// master.Run()
			} else {
				master.Run()
			}
		}

	}()
	for {
		time.Sleep(time.Second)
	}
}
