package elevator

import (
	"../def"
	"../driver"
	"fmt"
	"time"
)

var time50ms = 50 * time.Millisecond

func ExecuteOrders(orderListForExecuteOrders <-chan def.Orders, completedCurrentOrder chan<- bool, elevatorStateToMasterChan chan<- def.ElevatorState, elevatorStateChanForPrinting chan<- def.ElevatorState){

 	elevatorState := findDefinedElevatorState()
 	lastElevatorState := def.ElevatorState{}
 	orderList := def.Orders{}
 	floorSensorValue := -1

	for {
		select {
		case orderList = <-orderListForExecuteOrders:
			fmt.Println("new orderList In Execute orders. ", orderList)
			// newOrderList = true
			if len(orderList.Orders) > 0 {
				// fmt.Println("Executing order to Floor: ", orderList.Orders[0].Floor, ", with direction: ", orderList.Orders[0].Direction)

				// You're allready on desired floor
				if driver.Elev_get_floor_sensor_signal() == orderList.Orders[0].Floor {
					// driver.Elev_set_door_open_lamp(1)
					fmt.Println("You are allready on the desired floor")

				} else { /*You are not on the desired floor*/
					fmt.Println("You are not on your desired floor.")
					driver.Elev_set_door_open_lamp(0)
					if elevatorState.LastFloor == orderList.Orders[0].Floor {
						fmt.Println("elevatorState.lastFloor == destinationFloor")
						if orderList.Orders[0].Direction == 1 {
							driver.Elev_set_motor_direction(def.DIR_DOWN)
							elevatorState.Direction = def.DIR_DOWN
						} else {
							driver.Elev_set_motor_direction(def.DIR_UP)
							elevatorState.Direction = def.DIR_UP
						}
					} else if elevatorState.LastFloor < orderList.Orders[0].Floor {
						driver.Elev_set_motor_direction(def.DIR_UP)
							elevatorState.Direction = def.DIR_UP
					} else {
						driver.Elev_set_motor_direction(def.DIR_DOWN)
							elevatorState.Direction = def.DIR_DOWN
					}
				} /*Motor Direction set*/
			}

		default:
			if !isEqualElevatorState(lastElevatorState, elevatorState){
				elevatorStateToMasterChan <- elevatorState
				elevatorStateChanForPrinting <- elevatorState
				lastElevatorState = elevatorState
			}
			if len(orderList.Orders) > 0 {
								
				// Check if reached floor
				floorSensorValue = driver.Elev_get_floor_sensor_signal()
				if floorSensorValue >= 0 {
					elevatorState.LastFloor = floorSensorValue
					elevatorState.Destination = findDestination(orderList)
					driver.Elev_set_floor_indicator(floorSensorValue)
				}

				if floorSensorValue == orderList.Orders[0].Floor {
					fmt.Println("You reached your desired floor. Orderlist is now: ", orderList)
					completedCurrentOrder <- true
					time.Sleep(time.Millisecond * 150) //So the elevator stops in the middle of the sensor
					driver.Elev_set_motor_direction(def.DIR_STOP)
					driver.Elev_set_floor_indicator(orderList.Orders[0].Floor)
					driver.Elev_set_door_open_lamp(1)
					time.Sleep(1200 * time.Millisecond) // Keep door open
					driver.Elev_set_door_open_lamp(0)
				} else if driver.Elev_get_floor_sensor_signal() == 0 { /*This is just to be fail safe*/
					driver.Elev_set_motor_direction(def.DIR_UP)
					elevatorState.Direction = def.DIR_UP
				} else if driver.Elev_get_floor_sensor_signal() == def.N_FLOORS-1 {
					driver.Elev_set_motor_direction(def.DIR_DOWN)
					elevatorState.Direction = def.DIR_DOWN
				}
				time.Sleep(20 * time.Millisecond)
			}
		}
	}
}
func findDefinedElevatorState() def.ElevatorState {
	elevatorState := def.ElevatorState{Destination: def.IDLE}

	if driver.Elev_get_floor_sensor_signal() >= 0 {
		elevatorState.LastFloor = driver.Elev_get_floor_sensor_signal()
	}else{
		driver.Elev_set_motor_direction(def.DIR_UP)
		elevatorState.Direction = def.DIR_UP
		for {
			if driver.Elev_get_floor_sensor_signal() >= 0 {
				elevatorState.LastFloor = driver.Elev_get_floor_sensor_signal()
				driver.Elev_set_motor_direction(def.DIR_STOP)
				break
			}
			time.Sleep(10*time.Millisecond)
		} 
	}
	return elevatorState
}

func findDestination(orderList def.Orders) int {
	direction := orderList.Orders[0].Direction
	destination := orderList.Orders[0].Floor

	for _, order := range orderList.Orders{
		if order.Direction == direction {
			switch direction {
			case def.DIR_UP:
				if order.Floor > destination {
					destination = order.Floor
				}
			case def.DIR_DOWN:
				if order.Floor < destination {
					destination = order.Floor
				}
			}
		}else {
			return destination 
		}
	}
	// If every order is in the same direction
	return destination
}


func isEqualElevatorState(state1 def.ElevatorState, state2 def.ElevatorState) bool{
	return state1.LastFloor == state2.LastFloor && state1.Direction == state2.Direction && state1.Destination == state2.Destination
}

func CheckForElevatorFloorUpdates(updateElevatorStateFloor chan<- int) {
	tempLastFloor := -1
	for {
		lastFloor := driver.Elev_get_floor_sensor_signal()
		if lastFloor >= 0 && lastFloor < def.N_FLOORS && lastFloor != tempLastFloor {
			tempLastFloor = lastFloor
			if lastFloor == 0 {
				updateElevatorStateFloor <- 0
				driver.Elev_set_floor_indicator(lastFloor)
			} else if lastFloor < (def.N_FLOORS - 1) {
				driver.Elev_set_floor_indicator(lastFloor)
				updateElevatorStateFloor <- lastFloor
			} else if lastFloor == (def.N_FLOORS - 1) {
				driver.Elev_set_floor_indicator(lastFloor)
				updateElevatorStateFloor <- lastFloor
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
}
