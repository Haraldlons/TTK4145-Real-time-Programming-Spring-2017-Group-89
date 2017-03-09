package controller

import (
	"../definitions"
	// "../driver"
	// "./../controller"
	// "./src/network"
	"../buttons"
	//"./src/driver"
	"../storage"
	//"./src/master"
	"../elevator"

	//"./src/watchdog"
	// "elevator"
	// "network"
	// "storage"
	"fmt"
	"time"
	// "encoding/json"
)

var elevatorState = definitions.ElevatorState{2, 0}

// var orderList = definitions.orderList
func Run() {
	// Initializing
	storage.LoadElevatorStateFromFile(&elevatorState)
	storage.LoadOrderListFromFile(&orderList)
	go elevator.ExectueOrders()
	go watchdog.SendImAliveMessages()
	go watchdog.CheckForMasterAlive()
	go watchdog.CheckForUpdatesFromMaster()
	go watchdog.CheckForElevatorStateUpdates()

	buttonPressesUpdates := make(chan button)
	go checkForButtonPresses()

	// elevatorState := definitions.ElevatorState{2, 0}
	// fmt.Println("ElevatorState during initialization: ", elevatorState)
	// orderList := make(chan,
	go elevator.PrintLastFloorIfChanged(&elevatorState)

	// elevator.GoToFloor(3, &elevatorState)
	internalButtonsPressesChan := make(chan [definitions.N_FLOORS]int)
	externalButtonsPressesChan := make(chan [definitions.N_FLOORS][2]int)
	// orderList <- storage.GetOrdersFromFile(3)
	go buttons.Check_button_internal(internalButtonsPressesChan)
	go buttons.Check_button_external(externalButtonsPressesChan)
	// go elevator.ExecuteOrders(channel )

	// /*Make JSON object and send it*/
	// m := definitions.TestMessage{"Alice", "Hello", 1294706395881547000}
	// Enco
	// b, _ := json.Marshal(m)
	// fmt.Println("Json in byte:", b)
	// fmt.Println("length of bytearray: ", len(b))
	// fmt.Println("err2:", err2)
	// check(err2)

	// initialize()
	go printInternalPresses(internalButtonsPressesChan)
	go printExternalPresses(externalButtonsPressesChan)
	for {
		time.Sleep(time.Millisecond * 100)
	}
	return true
}

func Change_master() bool {
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
	stopCurrentOrder := make(chan int) // Doesn't matter which data type.
	isFirstButtonPress := true
	for {
		select {
		case <-internalButtonsChan:
			fmt.Println("Internal button pressed: ", <-internalButtonsChan)
			if !isFirstButtonPress {
				stopCurrentOrder <- 1
			} //Value in channel doesn't matter
			// if(saveOrderToFile) { go findFloorAndGoTo()}
			go findFloorAndGoTo(internalButtonsChan, <-internalButtonsChan, stopCurrentOrder)
			time.Sleep(time.Millisecond * 200)
			isFirstButtonPress = false
			// default:
			// fmt.Println("No button pressed")
		}
	}
}

func findFloorAndGoTo(kanal chan [definitions.N_FLOORS]int, buttonPresses [definitions.N_FLOORS]int, stopCurrentOrder chan int) {
	// fmt.Println("ButtonPresses: ", buttonPresses)
	array := buttonPresses
	for i := 0; i < definitions.N_FLOORS; i++ {
		if array[i] == 1 {
			// fmt.Println("Going to floorfdsf: ", i)
			elevator.GoToFloor(i, &elevatorState, stopCurrentOrder)
			// fmt.Println("goToFloor Ended", i, " ended")
		}
	}
}
