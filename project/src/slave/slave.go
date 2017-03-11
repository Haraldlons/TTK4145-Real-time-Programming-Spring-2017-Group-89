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

var elevatorState = definitions.ElevatorState{2, 0, 0}

// var orderList = definitions.orderList
func Run() {
	fmt.Println("Slave started!")
	// Initializing
	driver.Elev_init()
	driver.Elev_set_motor_direction(driver.DIRECTION_STOP)

	// Channel Definitions
	internalButtonsPressesChan := make(chan [definitions.N_FLOORS]int)
	externalButtonsPressesChan := make(chan [definitions.N_FLOORS][2]int)
	updateElevatorStateDirection := make(chan int)
	updateElevatorStateFloor := make(chan int)
	stopSendingImAliveMessage := make(chan bool)
	// stopRecievingImAliveMessage := make(chan bool)
	masterHasDiedChan := make(chan bool)

	updatedOrderList := make(chan int)

	///////////////////////////////////////////
	// Make manually orderList
	totalOrderList := definitions.Orders{}
	// orderList := definitions.Orders{definitions}
	// listOfNumbers := []int{0, 1, 2, 1, 3}
	// secondListOfNumbers := []int{-1, 1, 1, -1, 1}

	// for i := range listOfNumbers {
	// 	totalOrderList = definitions.Orders{append(totalOrderList.Orders, definitions.Order{Floor: listOfNumbers[i], Direction: secondListOfNumbers[i]})}
	// }
	// fmt.Println("printing totalOrderList:", totalOrderList)
	// storage.SaveOrdersToFile(1, totalOrderList)
	///////////////////////////////////////////

	storage.LoadOrdersFromFile(1, &totalOrderList)
	fmt.Println("Loaded totalOrderlist from a file. Result: ", totalOrderList)
	storage.LoadElevatorStateFromFile(&elevatorState)

	go elevator.CheckForElevatorFloorUpdates(&elevatorState, updateElevatorStateFloor)
	go elevator.ListenAfterElevatStateUpdatesAndSaveToFile(&elevatorState, updateElevatorStateDirection, updateElevatorStateFloor)

	go network.SendSlaveIsAliveRegularly(44, stopSendingImAliveMessage)
	go watchdog.CheckIfMasterIsAliveRegularly(masterHasDiedChan)

	go buttons.Check_button_internal(internalButtonsPressesChan)
	go buttons.Check_button_external(externalButtonsPressesChan)
	// go handleInternalButtonPresses(internalButtonsPressesChan)
	// go handleExternalButtonPresses(externalButtonsPressesChan)

	go printExternalPresses(externalButtonsPressesChan)
	go printInternalPresses(internalButtonsPressesChan)
	go elevator.ExecuteOrders(&totalOrderList, &elevatorState, updatedOrderList, updateElevatorStateDirection)

	// go network.RecieveJSON()

	// go checkForUpdatesFromMaster()
	// go sendUpdatesToMaster()

	// buttonPressesUpdates := make(chan button)
	// go checkForButtonPresses()

	// elevatorState := definitions.ElevatorState{2, 0}

	// go elevator.PrintLastFloorIfChanged(&elevatorState)
	updatedOrderList <- 1
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

func printExternalPresses(externalButtonsChan chan [definitions.N_FLOORS][2]int) {
	for {
		select {
		case <-externalButtonsChan:

			fmt.Println("\nExternal button pressed: ", <-externalButtonsChan)
			// go findFloorAngGoTo(externalButtonsChan)
			time.Sleep(time.Millisecond * 200)

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
