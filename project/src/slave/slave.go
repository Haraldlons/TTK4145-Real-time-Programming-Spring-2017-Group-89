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
	"sync"
)

// var elevatorState = definitions.ElevatorState{}

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

	mutex := &sync.Mutex{}

	// Channel Definitions
	internalButtonsPressesChan := make(chan [definitions.N_FLOORS]int)
	externalButtonsPressesChan := make(chan [definitions.N_FLOORS][2]int)
	stopSendingImAliveMessage := make(chan bool)
	// stopRecievingImAliveMessage := make(chan bool)
	masterHasDiedChan := make(chan bool)
	completedCurrentOrder := make(chan bool)
	orderListForExecuteOrders := make(chan definitions.Orders)
	updatedOrderList := make(chan definitions.Orders)
	orderListToExternalPresses := make(chan definitions.Orders)

	// Channels for listening to updates in elevatorState variables
	updateElevatorStateFloor := make(chan int)
	updateElevatorStateDirection := make(chan int)
	updateElevatorDestinationChan := make(chan int)
	updateElevatorState := make(chan definitions.ElevatorState)
	updateElevatorStateForUpdatesInOrderList := make(chan definitions.ElevatorState)

	// Channels for setting updates in elevatorState
	elevatorStateChanForExecuteOrders := make(chan definitions.ElevatorState)
	elevatorStateForExternalPresses := make(chan definitions.ElevatorState)

	// Channels for printing in a nice format
	elevatorStateChanForPrinting := make(chan definitions.ElevatorState)
	orderListChanForPrinting := make(chan definitions.Orders)
	lastSentMsgToMasterChanForPrinting := make(chan definitions.MSG_to_master)
	lastRecievedMSGFromMasterChanForPrinting := make(chan definitions.MSG_to_slave)

	// Channels for preparing msg to master
	orderListForSendingToMaster := make(chan definitions.Orders)
	elevatorStateToMaster := make(chan definitions.ElevatorState)
	extButToMaster := make(chan definitions.Order)
	sendMessageToMaster := make(chan bool)

	go printEntireElevatorNetworkOnUpdate(elevatorStateChanForPrinting, orderListChanForPrinting, lastSentMsgToMasterChanForPrinting, lastRecievedMSGFromMasterChanForPrinting, mutex)

	// elevatorStateChanMap := make(map[string] chan definitions.ElevatorState)
	// elevatorStateChanMap["forExecuteOrders"] = elevatorStateChanForExecuteOrders
	// elevatorStateChanMap["forExternalPresses"] = elevatorStateForExternalPresses
	// elevatorStateChanMap["forFloorUpdates"] = ElevatorStateForFloorUpdatesChan

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

	go storage.LoadElevatorStateFromFile(updateElevatorState)

	go elevator.CheckForElevatorFloorUpdates(updateElevatorStateFloor)
	// go elevator.ListenAfterElevatStateUpdatesAndSaveToFile(&elevatorState, updateElevatorStateDirection, updateElevatorStateFloor)

	go network.SendSlaveIsAliveRegularly(44, stopSendingImAliveMessage)
	go watchdog.CheckIfMasterIsAliveRegularly(masterHasDiedChan)

	go buttons.Check_button_internal(internalButtonsPressesChan)
	go buttons.Check_button_external(externalButtonsPressesChan)
	// go handleInternalButtonPresses(internalButtonsPressesChan)
	// go handleExternalButtonPresses(externalButtonsPressesChan)

	go printExternalPresses(externalButtonsPressesChan, orderListToExternalPresses, elevatorStateForExternalPresses, elevator_id, lastSentMsgToMasterChanForPrinting, extButToMaster, sendMessageToMaster)
	go printInternalPresses(internalButtonsPressesChan)

	go watchdog.TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList, orderListForExecuteOrders, completedCurrentOrder, orderListToExternalPresses, elevator_id, updateElevatorStateForUpdatesInOrderList, orderListChanForPrinting, lastSentMsgToMasterChanForPrinting, orderListForSendingToMaster)
	go network.ListenToMasterUpdates(updatedOrderList, elevator_id, lastRecievedMSGFromMasterChanForPrinting)
	go listenToUpdatesToElevatorStateAndSendOnChannels(updateElevatorState, elevatorStateChanForExecuteOrders, updateElevatorStateFloor, updateElevatorStateDirection, updateElevatorDestinationChan, updateElevatorStateForUpdatesInOrderList, elevatorStateForExternalPresses, elevatorStateChanForPrinting, elevatorStateToMaster)
	go elevator.ExecuteOrders(elevatorStateChanForExecuteOrders, orderListForExecuteOrders, updateElevatorStateDirection, completedCurrentOrder, updateElevatorDestinationChan)
	go sendUpdatesToMaster(elevator_id, elevatorStateToMaster, orderListForSendingToMaster, extButToMaster , sendMessageToMaster, lastSentMsgToMasterChanForPrinting )
	// fmt.Println("GOROUTINES HAVE STARTED!")
	// go watchdog.CheckIfElevatorIsStuck(executeOrdersIsAliveChan)

	for {
		select {
		case <-masterHasDiedChan:
			fmt.Println("Master is not alive.")
			stopSendingImAliveMessage <- true
			// stopRecievingImAliveMessage <- true
			time.Sleep(time.Second)
			return // Turns into master
		}
	}
}

func listenToUpdatesToElevatorStateAndSendOnChannels(updateElevatorState <-chan definitions.ElevatorState, elevatorStateChanForExecuteOrders chan<- definitions.ElevatorState, updateElevatorStateFloor <-chan int, updateElevatorStateDirection <-chan int, updateElevatorDestinationChan <-chan int, updateElevatorStateForUpdatesInOrderList chan<- definitions.ElevatorState, elevatorStateForExternalPresses chan<- definitions.ElevatorState, elevatorStateChanForPrinting chan<- definitions.ElevatorState, elevatorStateToMaster chan<- definitions.ElevatorState) {
	elevatorState := definitions.ElevatorState{}
	for {
		select {
		case elevatorState := <-updateElevatorState:
			fmt.Println("updateState")
			fmt.Println("--------------------------------")
			fmt.Println("Current updated elevator state:", elevatorState)
			fmt.Println("--------------------------------")
			fmt.Println("0")
			elevatorStateChanForExecuteOrders <- elevatorState
			fmt.Println("1")
			fmt.Println("2")
			updateElevatorStateForUpdatesInOrderList <- elevatorState
			fmt.Println("3")
			elevatorStateForExternalPresses <- elevatorState
			fmt.Println("4")
			elevatorStateChanForPrinting <- elevatorState
			fmt.Println("5")
			elevatorStateToMaster <- elevatorState
		case updatedFloor := <-updateElevatorStateFloor:
			fmt.Println("updateFloor")
			elevatorState.LastFloor = updatedFloor

			// fmt.Println("--------------------------------")
			// fmt.Println("Current updated elevator state:", elevatorState)
			// fmt.Println("--------------------------------")

			fmt.Println("6")
			elevatorStateChanForExecuteOrders <- elevatorState
			fmt.Println("7")
			updateElevatorStateForUpdatesInOrderList <- elevatorState
			fmt.Println("8")
			elevatorStateForExternalPresses <- elevatorState
			fmt.Println("9")
			elevatorStateChanForPrinting <- elevatorState
			fmt.Println("10")			
			fmt.Println("11")
			elevatorStateToMaster <- elevatorState
		case updateDirection := <-updateElevatorStateDirection:
			fmt.Println("12")
			elevatorState.Direction = updateDirection
			fmt.Println("13")

			// fmt.Println("--------------------------------")
			// fmt.Println("Current updated elevator state:", elevatorState)
			// fmt.Println("--------------------------------")

			elevatorStateChanForExecuteOrders <- elevatorState
			fmt.Println("14")
			fmt.Println("15")
			updateElevatorStateForUpdatesInOrderList <- elevatorState
			fmt.Println("16")
			elevatorStateForExternalPresses <- elevatorState
			fmt.Println("17")
			elevatorStateChanForPrinting <- elevatorState
			fmt.Println("18")
			elevatorStateToMaster <- elevatorState
			fmt.Println("19")
		case elevatorState.Destination = <-updateElevatorDestinationChan:
			fmt.Println("Mottok updateElevatorDestinationChan")
			fmt.Println("20")
			// fmt.Println("--------------------------------")
			// fmt.Println("Current updated elevator state:", elevatorState)
			// fmt.Println("--------------------------------")

			elevatorStateChanForExecuteOrders <- elevatorState
			fmt.Println("20")
			fmt.Println("21")
			updateElevatorStateForUpdatesInOrderList <- elevatorState
			fmt.Println("22")
			elevatorStateForExternalPresses <- elevatorState
			fmt.Println("23")
			elevatorStateChanForPrinting <- elevatorState
			fmt.Println("24")
			elevatorStateToMaster <- elevatorState
			fmt.Println("25")

		}
		time.Sleep(50*time.Millisecond)
	}
}



func sendUpdatesToMaster(elevator_id string, elevatorStateToMaster chan definitions.ElevatorState, orderListForSendingToMaster chan definitions.Orders, extButToMaster chan definitions.Order, sendMessageToMaster chan bool, lastSentMsgToMasterChanForPrinting chan definitions.MSG_to_master){

	elevatorState := definitions.ElevatorState{}
	orderList := definitions.Orders{}
	externalButtonpresses := []definitions.Order{}
	externalButtonpress := definitions.Order{}
	msg_to_master := definitions.MSG_to_master{Orders: orderList, ElevatorState: elevatorState, ExternalButtonPresses: externalButtonpresses, Id: elevator_id}

	for {
		select {
			case msg_to_master.Orders = <- orderListForSendingToMaster:
				fmt.Println("Updated orderListForSendingToMaster: ", msg_to_master.Orders)
			case msg_to_master.ElevatorState = <-elevatorStateToMaster:
				fmt.Println("Updated orderListForSendingToMaster: ", msg_to_master.ElevatorState)
			case externalButtonpress = <-extButToMaster:
				externalButtonpresses = append(externalButtonpresses, externalButtonpress)
				msg_to_master.ExternalButtonPresses = externalButtonpresses
				fmt.Println("Updated orderListForSendingToMaster: ", msg_to_master.ExternalButtonPresses)
			case <-sendMessageToMaster:
				fmt.Println("Message ready to be sent to master: ", msg_to_master)
				network.SendUpdatesToMaster(msg_to_master, lastSentMsgToMasterChanForPrinting)
				externalButtonpresses = []definitions.Order{}
		}
	}
}

func printEntireElevatorNetworkOnUpdate(elevatorStateChanForPrinting <-chan definitions.ElevatorState, orderListChanForPrinting <-chan definitions.Orders, lastSentMsgToMasterChanForPrinting <-chan definitions.MSG_to_master, lastRecievedMSGFromMasterChanForPrinting <-chan definitions.MSG_to_slave, mutex *sync.Mutex){
	elevatorState := definitions.ElevatorState{}
	orderList := definitions.Orders{}
	lastSentMsgToMaster := definitions.MSG_to_master{}
	lastRecievedMSGFromMaster := definitions.MSG_to_slave{}
	updateNrMap := make(map[string] int)
	updateNrMap["totalUpdates"] = 0
	updateNrMap["elevatorState"] = 0
	updateNrMap["orderList"] = 0
	updateNrMap["msgToMaster"] = 0
	updateNrMap["msgToSlave"] = 0
		for {
			select {
			case elevatorState = <- elevatorStateChanForPrinting:
				updateNrMap["elevatorState"]++

				mutex.Lock()
				printNicely(updateNrMap,elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
				mutex.Unlock()

			case orderList = <- orderListChanForPrinting:
				updateNrMap["orderList"]++

				mutex.Lock()
				printNicely(updateNrMap,elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
				mutex.Unlock()

			case lastSentMsgToMaster = <- lastSentMsgToMasterChanForPrinting:
				updateNrMap["msgToMaster"]++

				mutex.Lock()
				printNicely(updateNrMap,elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
				mutex.Unlock()

			case lastRecievedMSGFromMaster = <-lastRecievedMSGFromMasterChanForPrinting:
				updateNrMap["msgToSlave"]++

				mutex.Lock() 
				printNicely(updateNrMap,elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
				mutex.Unlock()
			}
			updateNrMap["totalUpdates"]++
		}
}

func printNicely(updateNrMap map[string] int, elevatorState definitions.ElevatorState, orderList definitions.Orders, lastSentMsgToMaster definitions.MSG_to_master, lastRecievedMSGFromMaster definitions.MSG_to_slave){
	direction := ""
	if(elevatorState.Direction == 1){
		direction = "Up"
	}else if elevatorState.Direction == -1{
		direction = "Down"
	}else {
		direction = "Uknown"
	}

	lengthOrderList := len(orderList.Orders)

	print("\033[H\033[2J") /*Clears the screen*/
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("----------------------- Update NR: %d ----------------------\n", updateNrMap["totalUpdates"])
	fmt.Println("------------------------------------------------------------\n")
	fmt.Printf("ElevStatNr: %d, OrderListNr: %d, MsgTOMasterNr: %d, MsgTOSlaveNr: %d\n", updateNrMap["elevatorState"], updateNrMap["orderList"], updateNrMap["msgToMaster"], updateNrMap["msgToSlave"])
	fmt.Println("")
	fmt.Printf("State: LastFloor: %d (0-3), Direction: %s, Destination %d (0-3)\n", elevatorState.LastFloor,direction ,elevatorState.Destination)
	fmt.Println("")
	fmt.Println("OrderList: ", orderList)
	for i, orders := range orderList.Orders {
		fmt.Printf("%d. Floor: %d - Direction: %d\n",i+1, orders.Floor, orders.Direction )
	}
	for i := 0; i < (5-lengthOrderList); i++ {
		fmt.Println("")
	}
	fmt.Println("")
	fmt.Println("Last message sent to Master: ", lastSentMsgToMaster)
	fmt.Println("")
	fmt.Println("Last message recieved from Master: ", lastRecievedMSGFromMaster)
	fmt.Println("")
	fmt.Println("------------------------------------------------------------")
	fmt.Println("------------------------------------------------------------")
	fmt.Println("------------------------------------------------------------")

}


func Change_master() bool { /*Do we need this?*/
	return true
}

func printExternalPresses(externalButtonsChan chan [definitions.N_FLOORS][2]int, orderListToExternalPresses chan definitions.Orders, elevatorStateForExternalPresses <-chan definitions.ElevatorState, elevator_id string, lastSentMsgToMasterChanForPrinting chan<- definitions.MSG_to_master, extButToMaster chan definitions.Order, sendMessageToMaster chan bool) {
	// orderList := definitions.Orders{}
	// msg_to_master := definitions.MSG_to_master{}
	// go func() {
	// 	for {
	// 		select {
	// 		case orderList = <-orderListToExternalPresses:
	// 			// fmt.Println("OrderList to ExternalPresses has been updated with: ", orderList)
	// 			msg_to_master.Orders = orderList
	// 			// fmt.Println("MsgTOMaster in printExternalPresses: ", msg_to_master)
	// 		}
	// 	}
	// }()

	// elevatorState := definitions.ElevatorState{}

	// go func() {
	// 	for {
	// 		select {
	// 		case elevatorState = <-elevatorStateForExternalPresses:
	// 		}
	// 	}
	// }()

	for {
		select {
		case externalButtonPressed := <-externalButtonsChan:
			// fmt.Println("\nExternal button pressed: ", externalButtonPressed)
			externalButtonOrder := getOrderFromOneExternalPress(externalButtonPressed)
			fmt.Println("External button pressed with order: ", externalButtonOrder)
			// fmt.Println("------------------------------------")
			// fmt.Println("------------------------------------")
			// fmt.Println("msg_to_master.ExternalButtonPresses", msg_to_master.ExternalButtonPresses)
			// msg_to_master.ExternalButtonPresses = append(msg_to_master.ExternalButtonPresses, externalButtonOrder)
			// fmt.Println("------------------------------------")
			// fmt.Println("------------------------------------")
			// fmt.Println("msg_to_master.ExternalButtonPresses", msg_to_master.ExternalButtonPresses)
			// network.SendUpdatesToMaster(msg_to_master, elevatorState, elevator_id, lastSentMsgToMasterChanForPrinting)
			// msg_to_master.ExternalButtonPresses = []definitions.Order{}

			// go findFloorAngGoTo(externalButtonsChan)
			extButToMaster <- externalButtonOrder
			time.Sleep(time.Millisecond * 100)
			sendMessageToMaster <- true
		// case elevatorState = <-elevatorStateForExternalPresses:
			// if elevatorState.LastFloor > 0 {}
		// case orderList = <-orderListToExternalPresses:
			// fmt.Println("Something", orderList)
		// fmt.Println("OrderList to ExternalPresses has been updated with: ", orderList)
		// msg_to_master.Orders = orderList


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
			// fmt.Println("Nedover! i etasje: ", i)
			return definitions.Order{Floor: i, Direction: -1}
		} else if externalButtonpressed[i][0] == 1 {
			// fmt.Println("Oppover! I etasje: ", i)
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
