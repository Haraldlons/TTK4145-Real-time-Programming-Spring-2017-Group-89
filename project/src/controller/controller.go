package controller

import (
	// "../definitions"
	// "../driver"
	// "./../controller"
	// "./src/network"
	"../buttons"
	//"./src/driver"
	// "../storage"
	//"./src/master"
	//"./src/watchdog"
	// "elevator"
	// "network"
	// "storage"
	// "fmt"
	// "time"
)

// Channels



// func initialize() bool {
// 	orderList := make(chan, 
// 	internalButtonsChan := make(chan, int)
// 	orderList <- storage.GetOrdersFromFile(3)
// 	return true
// }

func Run() bool {
	// Initializing
	// orderList := make(chan, 
	internalButtonsChan := make(chan int)
	// orderList <- storage.GetOrdersFromFile(3)
	go buttons.Check_button_external(internalButtonsChan)
	// go elevator.ExecuteOrders(channel )

	// initialize()
	printInternalPresses(internalButtonsChan)


	


	return true
}

func Change_master() bool {
	return true
}

func printInternalPresses(internalButtonsChan chan int) {
	for {
		select {
			case <-internalButtonsChan:
				// fmt.Println("Internal button pressed: ", <-internalButtonsChan)
				// time.Sleep(time.Millisecond*40)
				
			// default:
				// fmt.Println("No button pressed")
		}
	}
}
