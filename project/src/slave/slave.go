package slave

import (
	"../definitions"
	"../driver"
	// "./../slave"
	"../buttons"
	"../elevator"
	"../network"
	//"./src/driver"
	"../storage"
	//"./src/master"
	"../watchdog"
	// "network"
	// "storage"
	"fmt"
	"time"
	// "encoding/json"
)

var elevatorState = definitions.ElevatorState{}

// var orderList = definitions.orderList
func Run() {
	fmt.Println("Slave started!")
	// Initializing
	driver.Elev_init()
	// value := driver.Elev_init()
	// fmt.Println("Return value from Elev_init()", value)

	driver.Elev_set_motor_direction(driver.DIRECTION_STOP)

	//Get IP address
	elevator_id, err := network.GetLocalIP()
	if err != nil {
		elevator_id = "localhost"
	}

	// Channel Definitions
	internalButtonsPressesChan := make(chan [definitions.N_FLOORS]int)
	externalButtonsPressesChan := make(chan [definitions.N_FLOORS][2]int)
	updateElevatorStateDirection := make(chan int)
	updateElevatorStateFloor := make(chan int)
	stopSendingImAliveMessage := make(chan bool)
	// stopRecievingImAliveMessage := make(chan bool)
	masterHasDiedChan := make(chan bool)
	completedCurrentOrder := make(chan bool)
	orderListForExecuteOrders := make(chan definitions.Orders)
	updatedOrderList := make(chan definitions.Orders)
	orderListToExternalPresses := make(chan definitions.Orders)
	///////////////////////////////////////////
	// Make manually orderList
	totalOrderList := definitions.Orders{}
	// orderList := definitions.Orders{definitions}
	listOfNumbers := []int{0, 3}
	secondListOfNumbers := []int{1, -1}

	for i := range listOfNumbers {
		totalOrderList = definitions.Orders{append(totalOrderList.Orders, definitions.Order{Floor: listOfNumbers[i], Direction: secondListOfNumbers[i]})}
	}
	// fmt.Println("printing totalOrderList:", totalOrderList)
	// storage.SaveOrdersToFile(1, totalOrderList)
	///////////////////////////////////////////

	// storage.LoadOrdersFromFile(1, &totalOrderList)
	// fmt.Println("Loaded totalOrderlist from a file. Result: ", totalOrderList)Sending JSON over network. Interface:  {{[]} {1 0 0} [] 129.241.187.151}

	storage.LoadElevatorStateFromFile(&elevatorState)

	go elevator.CheckForElevatorFloorUpdates(&elevatorState, updateElevatorStateFloor)
	go elevator.ListenAfterElevatStateUpdatesAndSaveToFile(&elevatorState, updateElevatorStateDirection, updateElevatorStateFloor)

	go network.SendSlaveIsAliveRegularly(44, stopSendingImAliveMessage)
	go watchdog.CheckIfMasterIsAliveRegularly(masterHasDiedChan)

	go buttons.Check_button_internal(internalButtonsPressesChan)
	go buttons.Check_button_external(externalButtonsPressesChan)
	// go handleInternalButtonPresses(internalButtonsPressesChan)
	// go handleExternalButtonPresses(externalButtonsPressesChan)

	go printExternalPresses(externalButtonsPressesChan, orderListToExternalPresses)
	go printInternalPresses(internalButtonsPressesChan)
	go elevator.ExecuteOrders(&elevatorState, orderListForExecuteOrders, updateElevatorStateDirection, completedCurrentOrder)

	go watchdog.TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList, orderListForExecuteOrders, completedCurrentOrder, orderListToExternalPresses, elevator_id, elevatorState)

	go network.ListenToMasterUpdates(updatedOrderList, elevator_id)
	// go sendUpdatesToMaster()

	// buttonPressesUpdates := make(chan button)
	// go checkForButtonPresses()

	// elevatorState := definitions.ElevatorState{2, 0}

	// go elevator.PrintLastFloorIfChanged(&elevatorState)
	// updatedOrderList <- 1

	// dfasfdf
	// time.Sleep(2 * time.Second)
	// fmt.Println("Sending JSON TO MASTER")
	time.Sleep(time.Second)

	// externalButtonsPress := <-externalButtonsPressesChan,
	// externalButtonsPresses := []Order{externalButtonsPress}

	// msg := definitions.MSG_to_master{Orders: totalOrderList, ExternalButtonPresses: []definitions.Order{definitions.Order{Floor: 2, Direction: -1},
	// 	definitions.Order{Floor: 0, Direction: 1},
	// 	definitions.Order{Floor: 3, Direction: -1},
	// 	definitions.Order{Floor: 1, Direction: -1},
	// }}
	// // msg := definitions.MSG_to_master{Orders: totalOrderList, Id: elevator_id, ExternalButtonPresses: []definitions.Order{definitions.Order{Floor: 2, Direction: -1}}}

	// fmt.Println("Sending from slave:", elevator_id, ", Message: ", msg)
	// network.SendUpdatesToMaster(msg)

	// newOrderList := definitions.Orders{}
	// listOfNumbers := []int{0, 1, 2, 1, 3}
	// secondListOfNumbers := []int{-1, 1, 1, -1, 1}

	// for i := range listOfNumbers {
	// 	newOrderList = definitions.Orders{append(newOrderList.Orders, definitions.Order{Floor: listOfNumbers[i], Direction: secondListOfNumbers[i]})}
	// }

	// updatedOrderList <- newOrderList

	for {
		select {
		case <-masterHasDiedChan:
			fmt.Println("Master is not alive.")
			stopSendingImAliveMessage <- true
			// stopRecievingImAliveMessage <- true
			time.Sleep(time.Second)
			return
		}
		// updatedOrderList <- 1
		// fmt.Println("updatedOrderList now!")
	}
	return
}

func Change_master() bool { /*Do we need this?*/
	return true
}

func printExternalPresses(externalButtonsChan chan [definitions.N_FLOORS][2]int, orderListToExternalPresses chan definitions.Orders) {
	orderList := definitions.Orders{}
	msg_to_master := definitions.MSG_to_master{}
	go func() {
		for {
			select {
			case orderList = <-orderListToExternalPresses:
				// fmt.Println("OrderList to ExternalPresses has been updated with: ", orderList)
				msg_to_master.Orders = orderList
				// fmt.Println("MsgTOMaster in printExternalPresses: ", msg_to_master)
			}
		}
	}()

	for {
		select {
		case externalButtonPressed := <-externalButtonsChan:

			// fmt.Println("\nExternal button pressed: ", externalButtonPressed)
			externalButtonOrder := getOrderFromOneExternalPress(externalButtonPressed)
			// fmt.Println("------------------------------------")
			// fmt.Println("------------------------------------")
			// fmt.Println("msg_to_master.ExternalButtonPresses", msg_to_master.ExternalButtonPresses)
			msg_to_master.ExternalButtonPresses = append(msg_to_master.ExternalButtonPresses, externalButtonOrder)
			// fmt.Println("------------------------------------")
			// fmt.Println("------------------------------------")
			// fmt.Println("msg_to_master.ExternalButtonPresses", msg_to_master.ExternalButtonPresses)
			network.SendUpdatesToMaster(msg_to_master, elevatorState)
			msg_to_master.ExternalButtonPresses = []definitions.Order{}

			// go findFloorAngGoTo(externalButtonsChan)

			time.Sleep(time.Millisecond * 10)

			// default:
			// fmt.Println("No button pressed")
		}
	}
}

func printInternalPresses(internalButtonsChan chan [definitions.N_FLOORS]int) {
	// stopCurrentOrder := make(chan int) // Doesn't matter which data type.
	// isFirstButtonPress := true
	for {
		select {
		case list := <-internalButtonsChan:
			// print("\033[H\033[2J")
			fmt.Println("Internal button pressed: ", list)
			driver.Elev_set_button_lamp(2, getFloorFromInternalPress(list), 1)
			// if !isFirstButtonPress {
			// 	stopCurrentOrder <- 1
			// } //Value in channel doesn't matter
			// go findFloorAndGoTo(internalButtonsChan, <-internalButtonsChan, stopCurrentOrder)
			time.Sleep(time.Millisecond * 200)
			// isFirstButtonPress = false
			// default:
			// 	fmt.Println("No button pressed")
			// 	time.Sleep(time.Millisecond * 500)
			driver.Elev_set_button_lamp(2, getFloorFromInternalPress(list), 0)
		}
	}
}

func getOrderFromOneExternalPress(externalButtonpressed [definitions.N_FLOORS][2]int) definitions.Order {
	// fmt.Println("externalButtonpressed:", externalButtonpressed)
	for i := 0; i < definitions.N_FLOORS; i++ {
		if externalButtonpressed[i][1] == 1 {
			fmt.Println("Nedover! i etasje: ", i)
			return definitions.Order{Floor: i, Direction: -1}
		} else if externalButtonpressed[i][0] == 1 {
			fmt.Println("Oppover! I etasje: ", i)
			return definitions.Order{Floor: i, Direction: 1}
		}
	}
	return definitions.Order{}
}

func getFloorFromInternalPress(buttonPresses [definitions.N_FLOORS]int) int {
	array := buttonPresses
	for i := 0; i < definitions.N_FLOORS; i++ {
		if array[i] == 1 {
			// fmt.Println("Going to floorfdsf: ", i)
			return i
			// elevator.GoToFloor(i, &elevatorState, stopCurrentOrder)
			// fmt.Println("goToFloor Ended", i, " ended")
		}
	}
	return 0
}

// func findFloorAndGoTo(kanal chan [definitions.N_FLOORS]int, buttonPresses [definitions.N_FLOORS]int, stopCurrentOrder chan int) {
// 	// fmt.Println("ButtonPresses: ", buttonPresses)
// 	array := buttonPresses
// 	for i := 0; i < definitions.N_FLOORS; i++ {
// 		if array[i] == 1 {
// 			// fmt.Println("Going to floorfdsf: ", i)
// 			elevator.GoToFloor(i, &elevatorState, stopCurrentOrder)
// 			// fmt.Println("goToFloor Ended", i, " ended")
// 		}
// 	}
// }
