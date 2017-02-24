package main

import (
	//"master"
	"./src/driver"
	"./src/definitions"
	//"./src/buttons"
	//"./src/driver"
	// "./src/storage"
	//"./src/master"
	//"./src/watchdog"
	//"./src/network"
	"fmt"
	//"os"
	//"time"
)

func main() {
	fmt.Println("Main function started")
	//network.Run()

	driver.Elev_init();
	driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)

	elevatorState := definitions.ElevatorState{driver.Elev_get_floor_sensor_signal(),0}
	fmt.Println("elevatorInfo during initialization: ", elevatorState)

	stopSignal := 0 

	for true {
		// fmt.Println("Elev_get_floor_sensor_signal: ", driver.Elev_get_floor_sensor_signal())
		printLastFloorIfChanged(&elevatorState)
		// updateElevatorStateIfChanged(&elevatorState)

		//driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)

		if driver.Elev_get_floor_sensor_signal() == driver.N_FLOORS - 1 {
			// fmt.Println("Bobby Brown")
			driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
		} else if driver.Elev_get_floor_sensor_signal() == 0{
			// fmt.Println("Bobby Brown inverse")
			driver.Elev_set_motor_direction(driver.DIRECTION_UP)
		}
		
		stopSignal = driver.Elev_get_stop_signal()
		if stopSignal != 0 {
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			fmt.Println("Stopping program, with stop signal: ", stopSignal)
			fmt.Println("Another call to Elev_get_stop_signal(): ", driver.Elev_get_stop_signal())
			return
		}
	}

}

func printLastFloorIfChanged(elevatorState *definitions.ElevatorState) {
	lastFloor := driver.Elev_get_floor_sensor_signal();
	switch lastFloor {
				case 0: 
					if elevatorState.LastFloor != 0 {
						elevatorState.LastFloor = 0
						fmt.Println("Last Floor: 1. elevatorState: ", elevatorState)
					}
				case 1: 
					if elevatorState.LastFloor != 1 {
						elevatorState.LastFloor = 1
						fmt.Println("Last Floor: 2. elevatorState: ", elevatorState)
					}
				case 2: 
					if elevatorState.LastFloor != 2 {
						elevatorState.LastFloor = 2
						fmt.Println("Last Floor: 3. elevatorState: ", elevatorState)
					}
				case 3: 
					if elevatorState.LastFloor != 3 {
						elevatorState.LastFloor = 3
						fmt.Println("Last Floor: 4. elevatorState: ", elevatorState)
					}
				default:

			}
}
