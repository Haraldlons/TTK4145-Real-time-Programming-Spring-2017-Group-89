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
	"time"
)
var delay = 50 * time.Millisecond
var endProgram = false
var	elevatorActive = false


func main() {
	fmt.Println("Main function started")
	//network.Run()

	driver.Elev_init();
	driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)

	elevatorState := definitions.ElevatorState{driver.Elev_get_floor_sensor_signal(),0}
	fmt.Println("elevatorInfo during initialization: ", elevatorState)

	stopSignal := 0 
	// buttonSignal := driver.Elev_get_button_signal(0,0)

	driver.Elev_set_floor_indicator(3)
	goToFirstFloor := 0
	goToSecondFloor := 0
	goToThirdFloor := 0
	goToFourthFloor := 0

	go goToFloor(3, &elevatorState)

	for {
		// fmt.Println("Elev_get_floor_sensor_signal: ", driver.Elev_get_floor_sensor_signal())
		printLastFloorIfChanged(&elevatorState)
		// updateElevatorStateIfChanged(&elevatorState)

		//driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)

		// if driver.Elev_get_floor_sensor_signal() == driver.N_FLOORS - 1 {
		// 	// fmt.Println("Bobby Brown")
		// 	driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
		// } else if driver.Elev_get_floor_sensor_signal() == 0{
		// 	// fmt.Println("Bobby Brown inverse")
		// 	driver.Elev_set_motor_direction(driver.DIRECTION_UP)
		// }
		goToFirstFloor = driver.Elev_get_button_signal(2,0)
		goToSecondFloor = driver.Elev_get_button_signal(2,1)
		goToThirdFloor = driver.Elev_get_button_signal(2,2)
		goToFourthFloor = driver.Elev_get_button_signal(2,3)

		if(goToFirstFloor == 1) { go goToFloor(0, &elevatorState) }
		if(goToSecondFloor == 1) { go goToFloor(1, &elevatorState) }
		if(goToThirdFloor == 1) { go goToFloor(2, &elevatorState) }
		if(goToFourthFloor == 1) { go goToFloor(3, &elevatorState) }

		setFloorIndicator()

		// fmt.Println(goToFirstFloor, goToSecondFloor, goToThirdFloor, goToFourthFloor)

		// if(buttonSignal != 0){
		// 	go goToFloor(buttonSignal, &elevatorState)
		// 	fmt.Println("buttonSignal: ", buttonSignal)
		// 	time.Sleep(delay)
		// }else {
		// 	fmt.Println(buttonSignal)
		// 	time.Sleep(delay)
		// }

		if(endProgram){
			// driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			// fmt.Println("endProgram == true. Stopping program")
			return
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
						elevatorState.Direction = definitions.DIR_UP
						elevatorState.LastFloor = 0
						fmt.Println("Last Floor: 1. Direction: ", elevatorState.Direction)
					}
				case 1: 
					if elevatorState.LastFloor != 1 {
						if(elevatorState.LastFloor > lastFloor){
							elevatorState.Direction = definitions.DIR_DOWN
						}else {
							elevatorState.Direction = definitions.DIR_UP
						}
						elevatorState.LastFloor = 1
						fmt.Println("Last Floor: 2. Direction: ", elevatorState.Direction)
					}
				case 2: 
					if elevatorState.LastFloor != 2 {
						if(elevatorState.LastFloor > lastFloor){
							elevatorState.Direction = definitions.DIR_DOWN
						}else {
							elevatorState.Direction = definitions.DIR_UP
						}
						elevatorState.LastFloor = 2
						fmt.Println("Last Floor: 3. Direction: ", elevatorState.Direction)
					}
				case 3: 
					if elevatorState.LastFloor != 3 {
						elevatorState.Direction = definitions.DIR_DOWN
						elevatorState.LastFloor = 3
						fmt.Println("Last Floor: 4. Direction: ", elevatorState.Direction)
					}
				default:

			}
}

func goToFloor(destinationFloor int, elevatorState *definitions.ElevatorState) {
	if(!elevatorActive){
		elevatorActive = true

		fmt.Println("Going to floor: ", destinationFloor+1)
		// direction := elevatorState.Direction
		lastFloor := elevatorState.LastFloor

		if(driver.Elev_get_floor_sensor_signal() == destinationFloor){
			fmt.Println("You are allready on the desired floor")
			elevatorActive = false
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			// endProgram = true
			return
		}else {  /*You are not on the desired floor*/
				driver.Elev_set_door_open_lamp(0)
			if(lastFloor < destinationFloor){
				driver.Elev_set_motor_direction(driver.DIRECTION_UP)
			}else {
				driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
			}
			for {
				if(driver.Elev_get_floor_sensor_signal() == destinationFloor){
					fmt.Println("You reached your desired floor. Walk out\n")
					elevatorActive = false
					// driver.Elev_set_button_lamp(1,1,1)
					// driver.Elev_set_button_lamp(0,1,1)
					driver.Elev_set_floor_indicator(destinationFloor)
					driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
					// endProgram = true
					time.Sleep(delay*10)
					driver.Elev_set_door_open_lamp(1)
					return
				}else {
					time.Sleep(delay)
				}
			}
		}
	}
}

func setFloorIndicator(){
	sensorValue := driver.Elev_get_floor_sensor_signal()
	if(sensorValue != -1){
		driver.Elev_set_floor_indicator(sensorValue)
	}
}