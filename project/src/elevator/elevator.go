package elevator

import (
	"../def"
	"../driver"
	// "../storage"
	"fmt"
	"time"
)

var time50ms = 50 * time.Millisecond

func ExecuteOrders(orderListForExecuteOrders <-chan def.Orders, completedCurrentOrder chan<- bool, elevatorStateToMasterChan chan<- def.ElevatorState, elevatorStateChanForPrinting chan<- def.ElevatorState){

 	elevatorState := findDefinedElevatorState()
 	lastElevatorState := def.ElevatorState{}
 	orderList := def.Orders{}
 	floorSensorValue := -1

 	// time.Sleep(time.Millisecond*1000)



	// elevatorState := storage.LoadElevatorStateFromFile(updateElevatorState)
	
	

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
					// fmt.Println("You are not going places")
				} /*Motor Direction set*/
			}
		// case elevatorState = <-elevatorStateChanForExecuteOrders:

		default:
			if !isEqualElevatorState(lastElevatorState, elevatorState){
				elevatorStateToMasterChan <- elevatorState
				elevatorStateChanForPrinting <- elevatorState
				lastElevatorState = elevatorState
			}

			// fmt.Printf(".", len(orderList.Orders))
			if len(orderList.Orders) > 0 {
				// fmt.Println("driver.Elev_get_floor_sensor_signal() == orderList.Orders[0].Floor: ", driver.Elev_get_floor_sensor_signal(), ",", orderList.Orders[0].Floor)
								
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

					// fmt.Println("Finished with sending to channel updateElevatorStateDirection")
					driver.Elev_set_floor_indicator(orderList.Orders[0].Floor)
					driver.Elev_set_door_open_lamp(1)
					time.Sleep(500 * time.Millisecond) // Keep door open
					driver.Elev_set_door_open_lamp(0)
				} else if driver.Elev_get_floor_sensor_signal() == 0 { /*This is just to be fail safe*/
					driver.Elev_set_motor_direction(def.DIR_UP)
					// updateElevatorStateDirection <- def.DIR_UP
				} else if driver.Elev_get_floor_sensor_signal() == def.N_FLOORS-1 {
					driver.Elev_set_motor_direction(def.DIR_DOWN)
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

// func updateElevatorDestination(orderList def.Orders, elevatorState def.ElevatorState) int {
// 	lastCheckedFloorInOrderList := 0
// 	if len(orderList.Orders) == 0 {
// 		return -1
// 	} else {
// 		switch elevatorState.Direction {
// 		case def.DIR_UP:
// 			lastCheckedFloorInOrderList = 0
// 			maxFloor := 0
// 			for _, orders := range orderList.Orders {
// 				if lastCheckedFloorInOrderList > order.Floor {
// 					break
// 				}
// 				lastCheckedFloorInOrderList = order.Floor
// 			}
// 			updateElevatorDestinationChan <- maxFloor
// 			break
// 		case def.DIR_DOWN:
// 			fmt.Println("3. Step")
// 			lastCheckedFloorInOrderList = def.N_FLOORS
// 			minFloor := def.N_FLOORS
// 			for _, orders := range orderList.Orders {
// 				if orders.Floor < minFloor {
// 					minFloor = orders.Floor
// 				}
// 			}
// 			updateElevatorDestinationChan <- minFloor
// 		default:
// 			fmt.Println("4. Step")
// 			updateElevatorDestinationChan <- orderList.Orders[0].Floor
// 			fmt.Println("5. Step")
// 		}
// 	}
// 	return
// }


// func ListenAfterElevatStateUpdatesAndSaveToFile(elevatorState *def.ElevatorState, updateElevatorStateDirection chan int, updateElevatorStateFloor chan int) {
// 	for {
// 		select {
// 		case tempDirection := <-updateElevatorStateDirection:
// 			elevatorState.Direction = tempDirection
// 			storage.SaveElevatorStateToFile(elevatorState)
// 			time.Sleep(10 * time.Millisecond)
// 		case tempFloor := <-updateElevatorStateFloor:
// 			elevatorState.LastFloor = tempFloor
// 			storage.SaveElevatorStateToFile(elevatorState)
// 			time.Sleep(10 * time.Millisecond)
// 		}
// 	}
// }

func CheckForElevatorFloorUpdates(updateElevatorStateFloor chan<- int) {
	tempLastFloor := -1
	for {
		lastFloor := driver.Elev_get_floor_sensor_signal()
		if lastFloor >= 0 && lastFloor < def.N_FLOORS && lastFloor != tempLastFloor {
			tempLastFloor = lastFloor
			if lastFloor == 0 {
				updateElevatorStateFloor <- 0
				// elevatorState.LastFloor = 0
				// fmt.Println("Last Floor: 1. Direction: ", elevatorState.Direction, "(maybe need to use * )")
				// storage.SaveElevatorStateToFile(elevatorState)
				driver.Elev_set_floor_indicator(lastFloor)
			} else if lastFloor < (def.N_FLOORS - 1) {
				// if elevatorState.LastFloor > lastFloor {
				// 	elevatorState.Direction = def.DIR_DOWN
				// } else {
				// 	elevatorState.Direction = def.DIR_UP
				// }
				driver.Elev_set_floor_indicator(lastFloor)
				updateElevatorStateFloor <- lastFloor
				// fmt.Println("Last Floor: ", lastFloor, ". Direction: ", elevatorState.Direction)
				// storage.SaveElevatorStateToFile(elevatorState)
			} else if lastFloor == (def.N_FLOORS - 1) {
				// elevatorState.Direction = def.DIR_DOWN
				// elevatorState.LastFloor = lastFloor
				driver.Elev_set_floor_indicator(lastFloor)
				updateElevatorStateFloor <- lastFloor
				// fmt.Println("Last Floor: ", def.N_FLOORS, ". Direction: ", elevatorState.Direction)
				// storage.SaveElevatorStateToFile(elevatorState)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
}

/*This functions should be cleaned up. I have an ide how to do it*/
// func PrintLastFloorIfChanged(elevatorState *def.ElevatorState) {
// }

func GoToFloor(destinationFloor int, elevatorState def.ElevatorState, stopCurrentOrder chan bool, completedCurrentOrder chan<- bool, updateElevatorStateDirection chan<- int, goToFloorIsAlive chan<- bool) {
	// defer fmt.Println("Exeting goToFloor to floor: ", destinationFloor)
	// storage.SaveOrderToFile(destinationFloor)
	// elevatorActive = true

	// go func(){
	// 			fmt.Println("Floor: ", driver.Elev_get_floor_sensor_signal())
	// 		time.Sleep(time.Second)
	// }()

	fmt.Println("Going to floor: ", destinationFloor, " (0-3) ")
	direction := elevatorState.Direction
	lastFloor := elevatorState.LastFloor

	if driver.Elev_get_floor_sensor_signal() == destinationFloor {
		// storage.SaveOrderToFile(-1)
		fmt.Println("You are allready on the desired floor")
		// elevatorActive = false

		driver.Elev_set_motor_direction(def.DIR_STOP)
		completedCurrentOrder <- true
		goToFloorIsAlive <- true
		// endProgram = true
		for {

			select {
			case <-stopCurrentOrder:
				fmt.Println("Finially got message to stop going to floor, ", destinationFloor)
				return
			case <-time.After(500 * time.Millisecond):
				fmt.Println("Still have not got message to kill this order to floor: ", destinationFloor)
			}
		}
		return
	} else { /*You are not on the desired floor*/
		// fmt.Println("You are not on the desired floor")
		driver.Elev_set_door_open_lamp(0)
		if lastFloor == destinationFloor {
			fmt.Println("lastFloor == destinationFloor")
			if direction == 1 {
				driver.Elev_set_motor_direction(def.DIR_DOWN)
				updateElevatorStateDirection <- def.DIR_DOWN
			} else {
				driver.Elev_set_motor_direction(def.DIR_UP)
				updateElevatorStateDirection <- def.DIR_UP
			}
		} else if lastFloor < destinationFloor {
			driver.Elev_set_motor_direction(def.DIR_UP)
			updateElevatorStateDirection <- def.DIR_UP
		} else {
			driver.Elev_set_motor_direction(def.DIR_DOWN)
			updateElevatorStateDirection <- def.DIR_DOWN
		}
		for {
			select {
			case <-stopCurrentOrder:
				fmt.Println("stopCurrentOrder recieved. Stopping to floor: ", destinationFloor)
				return
			default:
				goToFloorIsAlive <- true

				// fmt.Println("Floor: ", driver.Elev_get_floor_sensor_signal())
				// fmt.Println("Testing")
				if driver.Elev_get_floor_sensor_signal() == destinationFloor {
					// orderList <- orderList[1:]
					fmt.Println("You reached your desired floor. Walk out\n")
					updateElevatorStateDirection <- def.DIR_STOP

					time.Sleep(time.Millisecond * 150) //So the elevator stops in the middle of the sensor
					// elevatorActive = false
					// driver.Elev_set_button_lamp(1,1,1)
					// driver.Elev_set_button_lamp(0,1,1)
					driver.Elev_set_floor_indicator(destinationFloor)
					driver.Elev_set_motor_direction(def.DIR_STOP)
					// endProgram = true
					time.Sleep(time50ms * 10)
					driver.Elev_set_door_open_lamp(1)
					// storage.SaveOrderToFile(-1)
					time.Sleep(time.Millisecond * 100)
					goToFloorIsAlive <- true
					completedCurrentOrder <- true
					for {
						select {
						case <-stopCurrentOrder:
							fmt.Println("Finially got message to stop going to floor, ", destinationFloor)
							return
						case <-time.After(500 * time.Millisecond):
							goToFloorIsAlive <- true
							fmt.Println("Still have not got message to kill this order to floor: ", destinationFloor)
						}
					}
					return
				} else if driver.Elev_get_floor_sensor_signal() == 0 { /*This is just to be fail safe*/
					driver.Elev_set_motor_direction(def.DIR_UP)
					updateElevatorStateDirection <- def.DIR_UP
				} else if driver.Elev_get_floor_sensor_signal() == 3 {
					driver.Elev_set_motor_direction(def.DIR_DOWN)
				} else {
					time.Sleep(time50ms) // 50ms
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

// func GetState() struct {
// 	return state
// }
