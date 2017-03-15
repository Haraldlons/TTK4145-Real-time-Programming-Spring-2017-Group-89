package slave

import (
	"../def"
	"../driver"
	// "./../slave"
	"../buttons"
	"../elevator"
	"../network"
	//"./src/driver"
	// "../storage"
	//"./src/master"
	"../watchdog"
	// "network"
	// "storage"
	"fmt"
	"time"
	// "encoding/json"
	"sync"
)

// var elevatorState = def.ElevatorState{}

// var orderList = def.orderList
func Run() {
	fmt.Println("Slave started!")
	// Initializing
	driver.Elev_init()
	// value := driver.Elev_init()
	// fmt.Println("Return value from Elev_init()", value)

	driver.Elev_set_motor_direction(def.DIR_STOP)

	//Get IP address
	slave_id, err := network.GetLocalIP()
	if err != nil {
		slave_id = "localhost"
	}

	mutex := &sync.Mutex{}

	// Channel Definitions
	internalButtonsPressesChan := make(chan [def.N_FLOORS]int)
	externalButtonsPressesChan := make(chan [def.N_FLOORS][2]int)
	stopSendingImAliveMessage := make(chan bool)
	newInternalButtonOrderChan := make(chan def.Order)

	// stopRecievingImAliveMessage := make(chan bool)
	masterHasDiedChan := make(chan bool)
	completedCurrentOrder := make(chan bool)
	orderListForExecuteOrders := make(chan definitions.Orders)
	updatedOrderList := make(chan definitions.Orders)
	orderListForLightsChan := make(chan definitions.Orders)

	// Channels for listening to updates in elevatorState variables

	// Channels for setting updates in elevatorState

	// Channels for printing in a nice format
	elevatorStateChanForPrinting := make(chan def.ElevatorState)
	orderListChanForPrinting := make(chan def.Orders)
	lastSentMsgToMasterChanForPrinting := make(chan def.MSG_to_master)
	lastRecievedMSGFromMasterChanForPrinting := make(chan def.MSG_to_slave)

	// Channels for preparing msg to master
	orderListForSendingToMaster := make(chan def.Orders)
	elevatorStateToMasterChan := make(chan def.ElevatorState)
	extButToMaster := make(chan def.Order)
	sendMessageToMaster := make(chan bool)
	go printEntireElevatorNetworkOnUpdate(elevatorStateChanForPrinting, orderListChanForPrinting, lastSentMsgToMasterChanForPrinting, lastRecievedMSGFromMasterChanForPrinting, mutex)

	// elevatorStateChanMap := make(map[string] chan def.ElevatorState)
	// elevatorStateChanMap["forExecuteOrders"] = elevatorStateChanForExecuteOrders
	// elevatorStateChanMap["forExternalPresses"] = elevatorStateForExternalPresses
	// elevatorStateChanMap["forFloorUpdates"] = ElevatorStateForFloorUpdatesChan

	go network.SendSlaveIsAliveRegularly(slave_id, stopSendingImAliveMessage)
	go watchdog.CheckIfMasterIsAliveRegularly(masterHasDiedChan)

	go buttons.Check_button_internal(internalButtonsPressesChan)
	go buttons.Check_button_external(externalButtonsPressesChan)
	// go handleInternalButtonPresses(internalButtonsPressesChan)
	// go handleExternalButtonPresses(externalButtonsPressesChan)

	go network.ListenToMasterUpdates(updatedOrderList, slave_id, lastRecievedMSGFromMasterChanForPrinting /*, mutex*/)
	go printExternalPresses(externalButtonsPressesChan, slave_id, lastSentMsgToMasterChanForPrinting, extButToMaster, sendMessageToMaster)
	go printInternalPresses(internalButtonsPressesChan, newInternalButtonOrderChan, sendMessageToMaster)

	go watchdog.TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList, orderListForExecuteOrders, completedCurrentOrder, slave_id, orderListChanForPrinting, lastSentMsgToMasterChanForPrinting, orderListForSendingToMaster, sendMessageToMaster, newInternalButtonOrderChan, orderListForLightsChan)

	// go listenToUpdatesToElevatorStateAndSendOnChannels(updateElevatorState, elevatorStateChanForExecuteOrders, updateElevatorStateFloor, updateElevatorStateDirection, updateElevatorDestinationChan, elevatorStateChanForPrinting, elevatorStateToMaster)

	go elevator.ExecuteOrders(orderListForExecuteOrders, completedCurrentOrder, elevatorStateToMasterChan, elevatorStateChanForPrinting)

	go keepTrackOfExLights(orderListForLightsChan)

	go sendUpdatesToMaster(slave_id, elevatorStateToMasterChan, orderListForSendingToMaster, extButToMaster, sendMessageToMaster, lastSentMsgToMasterChanForPrinting, newInternalButtonOrderChan)

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

func listenToUpdatesToElevatorStateAndSendOnChannels(updateElevatorState <-chan def.ElevatorState, elevatorStateChanForExecuteOrders chan<- def.ElevatorState, updateElevatorStateFloor <-chan int, updateElevatorStateDirection <-chan int, updateElevatorDestinationChan <-chan int, elevatorStateChanForPrinting chan<- def.ElevatorState, elevatorStateToMaster chan<- def.ElevatorState) {
	elevatorState := def.ElevatorState{}

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
			fmt.Println("3")
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
			select {
			case elevatorState.Direction = <-updateElevatorStateDirection:
				elevatorStateChanForExecuteOrders <- elevatorState
				fmt.Println("7")
				fmt.Println("8")
				fmt.Println("9")
				elevatorStateChanForPrinting <- elevatorState
				fmt.Println("10")
				elevatorStateToMaster <- elevatorState
				fmt.Println("11")
			default:
				elevatorStateChanForExecuteOrders <- elevatorState
				fmt.Println("7")
				fmt.Println("8")
				fmt.Println("9")
				elevatorStateChanForPrinting <- elevatorState
				fmt.Println("10")
				elevatorStateToMaster <- elevatorState
				fmt.Println("11")
			}
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
			fmt.Println("16")
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
			select {
			case elevatorState.Direction = <-updateElevatorStateDirection:
			default:
				elevatorStateChanForExecuteOrders <- elevatorState
				fmt.Println("20")
				fmt.Println("21")
				fmt.Println("22")
				fmt.Println("23")
				elevatorStateChanForPrinting <- elevatorState
				fmt.Println("24")
				elevatorStateToMaster <- elevatorState
				fmt.Println("25")

			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func sendUpdatesToMaster(slave_id string, elevatorStateToMaster chan def.ElevatorState, orderListForSendingToMaster chan def.Orders, extButToMaster chan def.Order, sendMessageToMaster chan bool, lastSentMsgToMasterChanForPrinting chan def.MSG_to_master, newInternalButtonOrderChan chan def.Order) {

	elevatorState := def.ElevatorState{}
	orderList := def.Orders{}
	externalButtonpresses := []def.Order{}
	externalButtonpress := def.Order{}
	msg_to_master := def.MSG_to_master{Orders: orderList, ElevatorState: elevatorState, ExternalButtonPresses: externalButtonpresses, Id: slave_id}
	newInternalButtonPress := def.Order{}

	for {
		select {
		case msg_to_master.Orders = <-orderListForSendingToMaster:
			fmt.Println("Updated orderListForSendingToMaster: ", msg_to_master.Orders)
		case msg_to_master.ElevatorState = <-elevatorStateToMaster:
			fmt.Println("Updated orderListForSendingToMaster: ", msg_to_master.ElevatorState)
		case externalButtonpress = <-extButToMaster:
			externalButtonpresses = append(externalButtonpresses, externalButtonpress)
			msg_to_master.ExternalButtonPresses = externalButtonpresses
			fmt.Println("Updated orderListForSendingToMaster: ", msg_to_master.ExternalButtonPresses)
		case newInternalButtonPress = <-newInternalButtonOrderChan:
			externalButtonpresses = append(externalButtonpresses, newInternalButtonPress)
			msg_to_master.ExternalButtonPresses = externalButtonpresses
			fmt.Println("Updated orderListForSendingToMaster: ", msg_to_master.ExternalButtonPresses)
		case <-sendMessageToMaster:
			fmt.Println("Message ready to be sent to master: ", msg_to_master)
			network.SendUpdatesToMaster(msg_to_master, lastSentMsgToMasterChanForPrinting)
			externalButtonpresses = []def.Order{}
		}
	}
}

func printEntireElevatorNetworkOnUpdate(elevatorStateChanForPrinting <-chan def.ElevatorState, orderListChanForPrinting <-chan def.Orders, lastSentMsgToMasterChanForPrinting <-chan def.MSG_to_master, lastRecievedMSGFromMasterChanForPrinting <-chan def.MSG_to_slave, mutex *sync.Mutex) {
	elevatorState := def.ElevatorState{}
	orderList := def.Orders{}
	lastSentMsgToMaster := def.MSG_to_master{}
	lastRecievedMSGFromMaster := def.MSG_to_slave{}
	updateNrMap := make(map[string]int)
	updateNrMap["totalUpdates"] = 0
	updateNrMap["elevatorState"] = 0
	updateNrMap["orderList"] = 0
	updateNrMap["msgToMaster"] = 0
	updateNrMap["msgToSlave"] = 0
	for {
		select {
		case elevatorState = <-elevatorStateChanForPrinting:
			updateNrMap["elevatorState"]++

			mutex.Lock()
			printNicely(updateNrMap, elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
			mutex.Unlock()

		case orderList = <-orderListChanForPrinting:
			updateNrMap["orderList"]++

			mutex.Lock()
			printNicely(updateNrMap, elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
			mutex.Unlock()

		case lastSentMsgToMaster = <-lastSentMsgToMasterChanForPrinting:
			updateNrMap["msgToMaster"]++

			mutex.Lock()
			printNicely(updateNrMap, elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
			mutex.Unlock()

		case lastRecievedMSGFromMaster = <-lastRecievedMSGFromMasterChanForPrinting:
			updateNrMap["msgToSlave"]++

			mutex.Lock()
			printNicely(updateNrMap, elevatorState, orderList, lastSentMsgToMaster, lastRecievedMSGFromMaster)
			mutex.Unlock()
		}
		updateNrMap["totalUpdates"]++
	}
}

func printNicely(updateNrMap map[string]int, elevatorState def.ElevatorState, orderList def.Orders, lastSentMsgToMaster def.MSG_to_master, lastRecievedMSGFromMaster def.MSG_to_slave) {
	direction := ""
	if elevatorState.Direction == 1 {
		direction = "Up"
	} else if elevatorState.Direction == -1 {
		direction = "Down"
	} else {
		direction = "Uknown"
	}

	lengthOrderList := len(orderList.Orders)

	print("\033[H\033[2J") /*Clears the screen*/
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("----------------------- Update NR: %d ----------------------\n", updateNrMap["totalUpdates"])
	fmt.Println("------------------------------------------------------------\n")
	fmt.Printf("ElevStatNr: %d, OrderListNr: %d, MsgTOMasterNr: %d, MsgTOSlaveNr: %d\n", updateNrMap["elevatorState"], updateNrMap["orderList"], updateNrMap["msgToMaster"], updateNrMap["msgToSlave"])
	fmt.Println("")
	fmt.Printf("State: LastFloor: %d (0-3), Direction: %s, Destination %d (0-3)\n", elevatorState.LastFloor, direction, elevatorState.Destination)
	fmt.Println("")
	fmt.Println("OrderList: ", orderList)
	for i, orders := range orderList.Orders {
		fmt.Printf("%d. Floor: %d - Direction: %d\n", i+1, orders.Floor, orders.Direction)
	}
	for i := 0; i < (5 - lengthOrderList); i++ {
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

func keepTrackOfLights(orderListForLightsChan chan definitions.Orders) {

	orderList := definitions.Orders{}
	internalLightsOn := []bool{false, false, false, false}
	externalLightsOn := []bool{false, false, false, false}
	fmt.Println("externalLightsOn:", externalLightsOn)

	for {
		select {
		case orderList = <-orderListForLightsChan:
			internalLightsOn = []bool{false, false, false, false}
			externalLightsOn = []bool{false, false, false, false}

			for i, order := range orderList {
				if order.Direction == 0 {
					/*Internal light*/
					// Elev_set_button_lamp(button int, floor int, value int) {
					driver.Elev_set_button_lamp(2, order.Floor, 1)
					internalLightsOn[order.Floor] = true
				} else {
					/*External light*/
					driver.Elev_set_button_lamp(0, order.Floor, 1)
					driver.Elev_set_button_lamp(1, order.Floor, 1)
					externalLightsOn[order.Floor] = true
				}
			}
			for i := 0; i < def.N_FLOORS; i++ {
				if internalLightsOn[i] == false {
					driver.Elev_set_button_lamp(2, i, 0)
				}
				if externalLightsOn[i] == false {
					driver.Elev_set_button_lamp(1, i, 0)
					driver.Elev_set_button_lamp(0, i, 0)
				}
			}
		}
	}

}

func Change_master() bool { /*Do we need this?*/
	return true
}

func printExternalPresses(externalButtonsChan chan [def.N_FLOORS][2]int, slave_id string, lastSentMsgToMasterChanForPrinting chan<- def.MSG_to_master, extButToMaster chan def.Order, sendMessageToMaster chan bool) {
	// orderList := def.Orders{}
	// msg_to_master := def.MSG_to_master{}
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

	// elevatorState := def.ElevatorState{}

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
			// network.SendUpdatesToMaster(msg_to_master, elevatorState, slave_id, lastSentMsgToMasterChanForPrinting)
			// msg_to_master.ExternalButtonPresses = []def.Order{}

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

func printInternalPresses(internalButtonsChan <-chan [def.N_FLOORS]int, newInternalButtonOrderChan chan<- def.Order, sendMessageToMaster chan<- bool) {
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
			floor := getFloorFromInternalPress(list)
			fmt.Println("Decoded to going to FLOOR INTERNAL PRESS:", floor)
			newInternalButtonOrderChan <- def.Order{Floor: floor, Direction: 0}
			sendMessageToMaster <- true
			// isFirstButtonPress = false
			// default:
			// 	fmt.Println("No button pressed")
			time.Sleep(time.Millisecond * 300)
			driver.Elev_set_button_lamp(2, getFloorFromInternalPress(list), 1)
			time.Sleep(time.Millisecond * 200)
			driver.Elev_set_button_lamp(2, getFloorFromInternalPress(list), 0)

		}
	}
}

func getOrderFromOneExternalPress(externalButtonpressed [def.N_FLOORS][2]int) def.Order {
	// fmt.Println("externalButtonpressed:", externalButtonpressed)
	for i := 0; i < def.N_FLOORS; i++ {
		if externalButtonpressed[i][1] == 1 {
			// fmt.Println("Nedover! i etasje: ", i)
			return def.Order{Floor: i, Direction: -1}
		} else if externalButtonpressed[i][0] == 1 {
			// fmt.Println("Oppover! I etasje: ", i)
			return def.Order{Floor: i, Direction: 1}
		}
	}
	return def.Order{}
}

func getFloorFromInternalPress(buttonPresses [def.N_FLOORS]int) int {
	array := buttonPresses
	for i := 0; i < def.N_FLOORS; i++ {
		if array[i] == 1 {
			// fmt.Println("Going to floorfdsf: ", i)
			return i
			// elevator.GoToFloor(i, &elevatorState, stopCurrentOrder)
			// fmt.Println("goToFloor Ended", i, " ended")
		}
	}
	return 0
}
