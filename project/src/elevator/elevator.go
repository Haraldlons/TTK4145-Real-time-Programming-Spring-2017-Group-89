package elevator

import (
	"../definitions"
	"../driver"
	"../storage"
	"fmt"
	"time"
)

var delay = 50 * time.Millisecond
var endProgram = false

// var elevatorActive = false
// var elevatorState = definitions.ElevatorState{2, 0}
var msg = make([]byte, 8)

// func ExecuteOrders(orders struct, elevatorState *definitions.ElevatorState) {
// 	for {
// 		select {
// 		case orderListChanged <- orderList:
// 			go goToFloor(orderListChanged[0].floor)
// 		}
// 	}
// }

func GoToFloor(destinationFloor int, elevatorState *definitions.ElevatorState, stopCurrentOrder chan int) {
	defer fmt.Println("Exeting goToFloor to floor: ", destinationFloor)
	// storage.SaveOrderToFile(destinationFloor)
	// elevatorActive = true

	fmt.Println("Going to floor: ", destinationFloor, " (0-3) ")
	direction := elevatorState.Direction
	lastFloor := elevatorState.LastFloor

	if driver.Elev_get_floor_sensor_signal() == destinationFloor {
		// storage.SaveOrderToFile(-1)
		fmt.Println("You are allready on the desired floor")
		// elevatorActive = false
		driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
		// endProgram = true
		return
	} else { /*You are not on the desired floor*/
		fmt.Println("You are not on the desired floor")
		driver.Elev_set_door_open_lamp(0)
		if lastFloor == destinationFloor {
			fmt.Println("lastFloor == destinationFloor")
			if direction == 1 {
				driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
			} else {
				driver.Elev_set_motor_direction(driver.DIRECTION_UP)
			}
		} else if lastFloor < destinationFloor {
			driver.Elev_set_motor_direction(driver.DIRECTION_UP)
		} else {
			driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
		}
		for {
			select {
				case <- stopCurrentOrder:
					fmt.Println("stopCurrentOrder recieved. Stopping to floor: ", destinationFloor)
					return
				default:
					fmt.Println("Floor: ", driver.Elev_get_floor_sensor_signal())
					// fmt.Println("Testing")
					if driver.Elev_get_floor_sensor_signal() == destinationFloor {
						// orderList <- orderList[1:]
						fmt.Println("You reached your desired floor. Walk out\n")

						time.Sleep(time.Millisecond * 150) //So the elevator stops in the middle of the sensor
						// elevatorActive = false
						// driver.Elev_set_button_lamp(1,1,1)
						// driver.Elev_set_button_lamp(0,1,1)
						driver.Elev_set_floor_indicator(destinationFloor)
						driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
						// endProgram = true
						time.Sleep(delay * 10)
						driver.Elev_set_door_open_lamp(1)
						// storage.SaveOrderToFile(-1)
							for {
								select {
								case <- stopCurrentOrder:
									fmt.Println("Finially got message to stop going to floor, ", destinationFloor)
									return
								case <- time.After(2000 * time.Millisecond):
									fmt.Println("Still have not got message to kill this order to floor: ", destinationFloor)
								}
							}
						return
					} else if driver.Elev_get_floor_sensor_signal() == 0 { /*This is just to be fail safe*/
						driver.Elev_set_motor_direction(driver.DIRECTION_UP)
					} else if driver.Elev_get_floor_sensor_signal() == 3 {
						driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
					} else {
						time.Sleep(delay) // 50ms
					}
			}
		}
	}

}

func setFloorIndicator() {
	sensorValue := driver.Elev_get_floor_sensor_signal()
	if sensorValue != -1 {
		driver.Elev_set_floor_indicator(sensorValue)
	}
}

/*This functions should be cleaned up. I have an ide how to do it*/
func PrintLastFloorIfChanged(elevatorState *definitions.ElevatorState) {
	for {
		lastFloor := driver.Elev_get_floor_sensor_signal()
		switch lastFloor {
		case 0:
			if elevatorState.LastFloor != 0 {
				elevatorState.Direction = definitions.DIR_UP
				elevatorState.LastFloor = 0
				fmt.Println("Last Floor: 1. Direction: ", elevatorState.Direction)
				storage.SaveElevatorStateToFile(elevatorState)
			}
		case 1:
			if elevatorState.LastFloor != 1 {
				if elevatorState.LastFloor > lastFloor {
					elevatorState.Direction = definitions.DIR_DOWN
				} else {
					elevatorState.Direction = definitions.DIR_UP
				}
				elevatorState.LastFloor = 1
				fmt.Println("Last Floor: 2. Direction: ", elevatorState.Direction)
				storage.SaveElevatorStateToFile(elevatorState)
			}
		case 2:
			if elevatorState.LastFloor != 2 {
				if elevatorState.LastFloor > lastFloor {
					elevatorState.Direction = definitions.DIR_DOWN
				} else {
					elevatorState.Direction = definitions.DIR_UP
				}

				elevatorState.LastFloor = 2
				fmt.Println("Last Floor: 3. Direction: ", elevatorState.Direction)
				storage.SaveElevatorStateToFile(elevatorState)
			}
		case 3:
			if elevatorState.LastFloor != 3 {
				elevatorState.Direction = definitions.DIR_DOWN
				elevatorState.LastFloor = 3
				fmt.Println("Last Floor: 4. Direction: ", elevatorState.Direction)
				storage.SaveElevatorStateToFile(elevatorState)
			}

		default:

		}
		time.Sleep(time.Millisecond * 10)
	}
}

// func GetState() struct {
// 	return state
// }
