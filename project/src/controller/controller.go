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
)

var elevatorState = definitions.ElevatorState{2, 0}

// Channels

// func initialize() bool {
// 	orderList := make(chan,
// 	internalButtonsChan := make(chan, int)
// 	orderList <- storage.GetOrdersFromFile(3)
// 	return true
// }

func Run() bool {
	// Initializing

	// elevatorState := definitions.ElevatorState{2, 0}
	storage.ReadElevatorStateFromFile(&elevatorState)
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

	// initialize()
	go printInternalPresses(internalButtonsPressesChan)
	printExternalPresses(externalButtonsPressesChan)
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
		case internalButtonPresses := <-internalButtonsChan:
			fmt.Println("Internal button pressed: ", internalButtonPresses)
			if(!isFirstButtonPress){ stopCurrentOrder <- 1 } //Value in channel doesn't matter
 			go findFloorAndGoTo(internalButtonsChan, internalButtonPresses, stopCurrentOrder)
			time.Sleep(time.Millisecond * 100)
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
