package main

import (
	"./src/controller"
	// "./src/definitions"
	"./src/driver"
	// "./src/elevator"
	"./src/network"
	//"./src/buttons"
	//"./src/driver"
	// "./src/storage"
	//"./src/master"
	//"./src/watchdog"
	"fmt"
	// "log"
	// "os"
	"time"
	// "fmt"
	// "os/exec"
)

var delay = 50 * time.Millisecond
var endProgram = false
var elevatorActive = false

// var elevatorState = definitions.ElevatorState{2, 0}
var msg = make([]byte, 8)

func main() {
	fmt.Println("Main function started")

	driver.Elev_init()
	driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
	// Change only on Harald Branch

	if network.checkIfMasterAlreadyExist() {
		slave.run()
	} else {
		master.run()
	}

	for {
		// elevator.PrintLastFloorIfChanged(&elevatorState)
		// updateElevatorStateIfChanged(&elevatorState)

		// if driver.Elev_get_floor_sensor_signal() == driver.N_FLOORS - 1 {
		// 	// fmt.Println("Bobby Brown")
		// 	driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
		// } else if driver.Elev_get_floor_sensor_signal() == 0{
		// 	// fmt.Println("Bobby Brown inverse")
		// 	driver.Elev_set_motor_direction(driver.DIRECTION_UP)
		// }

		// goToFirstFloor = driver.Elev_get_button_signal(2, 0)
		// goToSecondFloor = driver.Elev_get_button_signal(2, 1)
		// goToThirdFloor = driver.Elev_get_button_signal(2, 2)
		// goToFourthFloor = driver.Elev_get_button_signal(2, 3)

		// if goToFirstFloor == 1 {
		// 	go goToFloor(0, &elevatorState)
		// 	// setOrderOverNetwork(0)
		// }
		// if goToSecondFloor == 1 {
		// 	go goToFloor(1, &elevatorState)
		// 	// setOrderOverNetwork(1)
		// }
		// if goToThirdFloor == 1 {
		// 	go goToFloor(2, &elevatorState)
		// 	// setOrderOverNetwork(2)
		// }
		// if goToFourthFloor == 1 {
		// 	go goToFloor(3, &elevatorState)
		// 	// setOrderOverNetwork(3)
		// }

		if endProgram {
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			fmt.Println("endProgram == true. Stopping program")
			return
		}

		stopSignal = driver.Elev_get_stop_signal()
		if stopSignal != 0 {
			// setOrderOverNetwork(0)
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			fmt.Println("Stopping program, with stop signal: ", stopSignal)
			fmt.Println("Another call to Elev_get_stop_signal(): ", driver.Elev_get_stop_signal())
			return
		}
	}
} //End main
