package watchdog

import (
	"../definitions"
	"../master"
	"../network"
	"../storage"
	// "net"
	"fmt"
	"time"
)

// var timeLimit = 10 * time.Second
// var listenTimer = 100 * time.Millisecond

// func SendNetworkAlive(udpBroadcast *net.UDPConn) bool { // Possibly unnesces
// 	msg := make([]byte, 1024)
// 	numLines, err := udpBroadcast.Write(msg)
// 	return err != nil
// }

func CheckIfMasterIsAliveRegularly(masterHasDiedChan chan bool) {
	masterIsAliveChan := make(chan string, 1)

	stopListening := make(chan bool)

	go network.ListenAfterAliveMasterRegularly(masterIsAliveChan, stopListening)
	master_id := ""
	for {
		select {
		case master_id = <-masterIsAliveChan:
			// fmt.Println("Master is still alive: , ", master_id)
		case <-time.After(time.Millisecond * 3000):
			fmt.Println("Master is not alive for the last three seconds")
			stopListening <- true
			fmt.Println("Has send stopListening signal to network.ListenAfterAliveMasterRegularly")
			masterHasDiedChan <- true
			return
		}
	}
	if master_id == "" {
	}
}

func TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList <-chan definitions.Orders, orderListForExecuteOrders chan<- definitions.Orders, completedCurrentOrder <-chan bool, elevator_id string, updateElevatorStateForUpdatesInOrderList <-chan definitions.ElevatorState, orderListChanForPrinting chan<- definitions.Orders, lastSentMsgToMasterChanForPrinting chan<- definitions.MSG_to_master, orderListForSendingToMaster chan definitions.Orders, internalPressOrderChan <-chan definitions.Order) {
	elevatorState := definitions.ElevatorState{}

	// go func() {
	// 	for {
	// 		select {
	// 		case elevatorState = <-updateElevatorStateForUpdatesInOrderList:
	// 		}
	// 	}
	// }()

	currentOrderList := definitions.Orders{}
	storage.LoadOrdersFromFile(1, &currentOrderList)
	fmt.Println("Loaded totalOrderlist from a file. Result: ", currentOrderList)
	orderListForExecuteOrders <- currentOrderList
	time.Sleep(50 * time.Millisecond)

	internalPressOrder := definitions.Order{}

	// time.Sleep(50 * time.Millisecond)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-completedCurrentOrder:
	// 			if len(currentOrderList.Orders) > 0 {
	// 				currentOrderList = definitions.Orders{currentOrderList.Orders[1:]}
	// 			}
	// 			orderListForExecuteOrders <- currentOrderList
	// 			orderListToExternalPresses <- currentOrderList
	// 			storage.SaveOrdersToFile(1, currentOrderList)
	// 			orderListChanForPrinting <- currentOrderList
	// 			msg := definitions.MSG_to_master{Orders: currentOrderList, Id: elevator_id}
	// 			// fmt.Println("msg_to_master: ", msg)
	// 			network.SendUpdatesToMaster(msg, elevatorState, elevator_id, lastSentMsgToMasterChanForPrinting)
	// 			time.Sleep(50 * time.Millisecond)
	// 		}
	// 	}
	// }()
	lastOrderList := definitions.Orders{}

	for {
		select {
		case currentOrderList = <-updatedOrderList:
			fmt.Println("Is this orderList going to currentOrderlist? ", currentOrderList, " lastOrderLIST:", lastOrderList)
			if checkIfChangedOrderList(lastOrderList, currentOrderList) {
				lastOrderList = currentOrderList
				time.Sleep(100 * time.Millisecond)
				fmt.Println("New Update to OrderList: ", currentOrderList)
				// if currentOrderList != updatedOrderList { /*If orderlist from master is identical to our copy*/
				fmt.Println("40")
				select {
				case <-completedCurrentOrder:
					fmt.Println("40,5")
				default:
					fmt.Println("41")
					orderListForExecuteOrders <- currentOrderList
					fmt.Println("42")
					orderListChanForPrinting <- currentOrderList
					fmt.Println("43")
					storage.SaveOrdersToFile(1, currentOrderList)
					fmt.Println("44")
					// sendOrderListUpdateToMaster(currentOrderList)
				}
				orderListForSendingToMaster <- currentOrderList
				fmt.Println("44,5")
			} else {
				fmt.Println("THEY ARE THE SAMMMEEMEMEMEMEMEMEMMEMEME")
			}

			time.Sleep(50 * time.Millisecond)
		case <-completedCurrentOrder:
			fmt.Printf("completedCurrentOrder")
			fmt.Println("CurrentOrderlist in special case:", currentOrderList)
			if len(currentOrderList.Orders) > 0 {
				fmt.Println("completedCurrentOrder23")
				currentOrderList = definitions.Orders{currentOrderList.Orders[1:]}
				fmt.Println("orderListafterSlice: ", currentOrderList)
			}
			orderListForExecuteOrders <- currentOrderList
			fmt.Println("45")
			fmt.Println("46")
			storage.SaveOrdersToFile(1, currentOrderList)
			orderListChanForPrinting <- currentOrderList
			fmt.Println("47")
			msg := definitions.MSG_to_master{Orders: currentOrderList, Id: elevator_id}
			// fmt.Println("msg_to_master: ", msg)
			fmt.Println("48")
			network.SendUpdatesToMaster(msg, lastSentMsgToMasterChanForPrinting)
			fmt.Println("49")
			orderListForSendingToMaster <- currentOrderList
			fmt.Println("50")
			time.Sleep(50 * time.Millisecond)
		case elevatorState = <-updateElevatorStateForUpdatesInOrderList:
			if elevatorState.LastFloor > 0 {
			}
		case internalPressOrder = <-internalPressOrderChan:
			currentOrderList = distributeInternalOrderToOrderList(internalPressOrder, currentOrderList, elevatorState)
			// }
		}
	}
}

func distributeInternalOrderToOrderList(internalPressOrder definitions.Order, currentOrderList definitions.Orders, elevatorState definitions.ElevatorState) definitions.Orders {
	newOrderList := definitions.Orders{}

	if master.CheckForDublicateOrder(currentOrderList, internalPressOrder.Floor) {
		return currentOrderList
	}

	tempNum := 0

	if elevatorState.LastFloor < internalPressOrder.Floor && internalPressOrder.Floor < elevatorState.Destination {
		// You are going up
		if currentOrderList.Orders[0] == elevatorState.Destination { /* You can add in front of currentOrderList */
			newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
			copy(newOrderList.Orders[1:], newOrderList.Orders[:])
			newOrderList.Orders[0] = internalPressOrder
			return newOrderList
		}else { /* There are orders before destinationOrder */
			for i, order := range currentOrderList.Orders{
				if order.Floor > tempNum {
					if order.Floor > internalPressOrder.Floor {
						// newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
						newOrderList.Orders = append(currentOrderList[i:], 
						copy(newOrderList.Orders[i+1:],  newOrderList.Orders[i:])
						newOrderList.Orders[0] = internalPressOrder
						return newOrderList
					}
					tempNum = order.Floor
				}else {

				}
			}

		}

	}else if {
		// You are going down
		tempNum = definitions.N_FLOOR-1
	}

	// if internalPressOrder.Floor > elevatorState.LastFloor && internalPressOrder.Direction == elevatorState.Direction {

	// }
}

func CheckIfElevatorIsStuck(executeOrdersIsAliveChan <-chan bool) {

	for {
		select {
		case <-executeOrdersIsAliveChan:
			fmt.Println("executeOrders Is still alive! Yeah")
		case <-time.After(time.Millisecond * 1000):
			fmt.Println("Have not recieved signal from executeOrders for 1 second")
		}
		time.Sleep(40 * time.Millisecond)
	}
}

func checkIfChangedOrderList(lastOrderList definitions.Orders, currentOrderList definitions.Orders) bool {
	if len(lastOrderList.Orders) != len(currentOrderList.Orders) {
		return true
	}
	for i, order := range lastOrderList.Orders {
		if order.Floor != currentOrderList.Orders[i].Floor {
			return true
		}
	}
	return false
}

/*
func CheckNetworkAlive(udpListen *net.UDPConn) int {
	listenChan := make(chan int, 1)
	lifeCheck := 1

	// Run listening goroutine
	go listen(listenChan, udpListen)

	for {
		select {
		case lifeCheck = <-listenChan:
			if lifeCheck == 1 {
				time.Sleep(listenTimer)
			} else { // Possibly dangerous. RETHINK!
				return -1
			}
		case <-time.After(timeLimit): // Node assumed dead
			return -1
		}
	}
}

func listen(listenChan chan int, udpListen *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		udpListen.ReadFromUDP(buf)
		listenChan <- int(buf)
		time.Sleep(listenTimer)
	}
}*/

/*
func CheckElevatorState(state var) var {
	state := true
	return state
}

func reset_Master() {

}

func reset_Elevator() {

}
*/
