package elevator

import (
	"fmt"
	"driver"
)

func ExecuteOrders(orders struct, elevatorState *definitions.ElevatorState) {


	for {
		select {
		case orderListChanged <- orderList:
			go goToFloor(orderListChanged[0].floor)
		}
	}
}

func goToFloor(destinationFloor int, elevatorState *definitions.ElevatorState) {
	if !elevatorActive {
		saveOrderToFile(destinationFloor)
		elevatorActive = true

		fmt.Println("Going to floor: ", destinationFloor+1)
		direction := elevatorState.Direction
		lastFloor := elevatorState.LastFloor

		if driver.Elev_get_floor_sensor_signal() == destinationFloor {
			saveOrderToFile(-1)

			fmt.Println("You are allready on the desired floor")
			elevatorActive = false
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			// endProgram = true
			return
		} else { /*You are not on the desired floor*/


			driver.Elev_set_door_open_lamp(0)
			if lastFloor == destinationFloor {
				if direction == 1 {
					driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
				} else {
					driver.Elev_set_motor_direction(driver.DIRECTION_UP)
				}
			}
			if lastFloor < destinationFloor {
				driver.Elev_set_motor_direction(driver.DIRECTION_UP)
			} else {
				driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
			}
			for {
				select {
				case <- orderList: // When orderList changes, we return the function
					return
				default:
					if driver.Elev_get_floor_sensor_signal() == destinationFloor {
						orderList <- orderList[1:]

						time.Sleep(time.Millisecond * 150) //So the elevator stops in the middle of the sensor
						fmt.Println("You reached your desired floor. Walk out\n")
						elevatorActive = false
						// driver.Elev_set_button_lamp(1,1,1)
						// driver.Elev_set_button_lamp(0,1,1)
						driver.Elev_set_floor_indicator(destinationFloor)
						driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
						// endProgram = true
						time.Sleep(delay * 10)
						driver.Elev_set_door_open_lamp(1)
						saveOrderToFile(-1)
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
}


func setFloorIndicator() {
	sensorValue := driver.Elev_get_floor_sensor_signal()
	if sensorValue != -1 {
		driver.Elev_set_floor_indicator(sensorValue)
	}
}

/*This functions should be cleaned up. I have an ide how to do it*/
func PrintLastFloorIfChanged(elevatorState *definitions.ElevatorState) {
	lastFloor := driver.Elev_get_floor_sensor_signal()
	switch lastFloor {
	case 0:
		if elevatorState.LastFloor != 0 {
			elevatorState.Direction = definitions.DIR_UP
			elevatorState.LastFloor = 0
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
			fmt.Println("Last Floor: 1. Direction: ", elevatorState.Direction)
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
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
		}
	case 2:
		if elevatorState.LastFloor != 2 {
			if elevatorState.LastFloor > lastFloor {
				elevatorState.Direction = definitions.DIR_DOWN
			} else {
				elevatorState.Direction = definitions.DIR_UP
			}

			elevatorState.LastFloor = 2
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
			fmt.Println("Last Floor: 3. Direction: ", elevatorState.Direction)
		}
	case 3:
		if elevatorState.LastFloor != 3 {
			elevatorState.Direction = definitions.DIR_DOWN
			elevatorState.LastFloor = 3
			fmt.Println("Last Floor: 4. Direction: ", elevatorState.Direction)
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
		}

	default:

	}
}

func GetState() struct {
	return state
}
