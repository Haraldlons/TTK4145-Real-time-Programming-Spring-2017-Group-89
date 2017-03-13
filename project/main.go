package main

import (
	"./src/slave"
	// "./src/definitions"
	//"./src/driver"
	// "./src/elevator"
	"./src/network"
	//"./src/buttons"
	"./src/driver"
	// "./src/storage"
	"./src/master"
	// "./src/watchdog"
	"fmt"
	// "log"
	// "os"
	// "os/exec"
	// "net"
	"time"
)

var delay = 50 * time.Millisecond
var elevatorActive = false
var port string = ":46723"

// var elevatorState = definitions.ElevatorState{2, 0}
var msg = make([]byte, 8)

func main() {
	fmt.Println("Main function started")
	// go slave.Run()
	// go master.Run()
	// go network.SetupNetwork()

	go func() {
		stopSignal := 0
		// time.Sleep(500 * time.Millisecond)
		for {
			stopSignal = driver.Elev_get_stop_signal()
			// fmt.Println("Stopsignal: ", stopSignal)
			if stopSignal != 0 {
				// setOrderOverNetwork(0)
				driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
				fmt.Println("Stopping program, with stop signal: ", stopSignal)
				fmt.Println("Another call to Elev_get_stop_signal(): ", driver.Elev_get_stop_signal())
				time.Sleep(100 * time.Millisecond)
				// return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	// udpAddr, _ := net.ResolveUDPAddr("udp", port)
	// fmt.Println("udpAddr", udpAddr)
	go func(){
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