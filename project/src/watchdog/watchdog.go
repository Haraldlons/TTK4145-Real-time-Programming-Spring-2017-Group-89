package watchdog

import (
	"../definitions"
	// "../master"
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
	masterIsAliveChan := make(chan string)

	stopListening := make(chan bool)

	go network.ListenAfterAliveMasterRegularly(masterIsAliveChan, stopListening)
	for {
		select {
		// case master_id := <-masterIsAliveChan:

		// fmt.Println("Master is still alive: , ", master_id)

		case <-time.After(time.Second * 3):
			fmt.Println("Master is not alive for the last three seconds")
			stopListening <- true
			fmt.Println("Has send stopListening signal to network.ListenAfterAliveMasterRegularly")
			masterHasDiedChan <- true
			return
		}
	}

	// if master_id == "" {
	// }
}

func TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList <-chan definitions.Orders, orderListForExecuteOrders chan<- definitions.Orders, completedCurrentOrder <-chan bool, elevator_id string, orderListChanForPrinting chan<- definitions.Orders, lastSentMsgToMasterChanForPrinting chan<- definitions.MSG_to_master, orderListForSendingToMaster chan definitions.Orders, newInternalButtonOrderChan chan definitions.Order, orderListForLightsChan chan<- definitions.Orders)

	currentOrderList := definitions.Orders{}
	storage.LoadOrdersFromFile(1, &currentOrderList)
	fmt.Println("Loaded totalOrderlist from a file. Result: ", currentOrderList)
	orderListForExecuteOrders <- currentOrderList
	time.Sleep(50 * time.Millisecond)

	lastOrderList := definitions.Orders{}

	for {
		select {
		case currentOrderList = <-updatedOrderList:
			fmt.Println("Is this orderList going to currentOrderlist? ", currentOrderList, " lastOrderLIST:", lastOrderList)
			if checkIfChangedOrderList(lastOrderList, currentOrderList) {
				lastOrderList = currentOrderList
				// time.Sleep(100 * time.Millisecond)
				fmt.Println("New Update to OrderList: ", currentOrderList)
				// if currentOrderList != updatedOrderList { /*If orderlist from master is identical to our copy*/
				fmt.Println("40")
				select {
				case <-completedCurrentOrder: /*This is a trick to avoid circular dependencies*/
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
				orderListForLightsChan <- currentOrderList
			} else {
				fmt.Println("THEY ARE THE SAMMMEEMEMEMEMEMEMEMMEMEME")
			}

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
			orderListForLightsChan<- currentOrderList
			fmt.Println("50,25")
		case newInternalButtonPress = <-newInternalButtonOrderChan:
			fmt.Println("50,5")
			currentOrderList = distributeInternalOrderToOrderList(newInternalButtonpress, currentOrderList, elevatorState)
			fmt.Println("51")
			orderListForExecuteOrders <- currentOrderList
			fmt.Println("52")
			orderListChanForPrinting <- currentOrderList
			fmt.Println("53")
			orderListForSendingToMaster <- currentOrderList
			fmt.Println("54")
			orderListForLightsChan <- currentOrderList
			fmt.Println("55")
		}
	}
}

func distributeInternalOrderToOrderList(internalPressOrder definitions.Order, currentOrderList definitions.Orders, elevatorState definitions.ElevatorState) definitions.Orders {
	newOrderList := definitions.Orders{}

	if master.CheckForDuplicateOrder(currentOrderList, internalPressOrder.Floor) {
		return currentOrderList
	}

	tempNum := 0

	if elevatorState.Direction == 1 {
		// You are going up
		fmt.Println("You are going up")
		if currentOrderList.Orders[0].Floor == elevatorState.Destination { /* You can add in front of currentOrderList */
			fmt.Println("First order is the destination floor")
			newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
			copy(newOrderList.Orders[1:], newOrderList.Orders[:])
			newOrderList.Orders[0] = internalPressOrder
			return newOrderList
		} else { /* There are orders before destinationOrder */
			for i, order := range currentOrderList.Orders {
				fmt.Println("Order[", i, "]: ", order)
				if order.Floor > tempNum { // To check where you turn
					fmt.Println(" order.Floor > tempNum ")
					if order.Floor > internalPressOrder.Floor && elevatorState.LastFloor < internalPressOrder.Floor {
						fmt.Println("This IF STATEMENT")
						newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
						copy(newOrderList.Orders[i+1:], newOrderList.Orders[i:])
						newOrderList.Orders[i] = internalPressOrder
						return newOrderList
					}
					tempNum = order.Floor
				}
				if tempNum == elevatorState.Destination {
					fmt.Println("tempNum == elevatorState.Destination")

					for j, order2 := range currentOrderList.Orders {
						fmt.Println("Length, ", len(currentOrderList.Orders), ", j, ", j)
						if j > i {
							if order2.Floor < internalPressOrder.Floor {
								fmt.Println("The other IF STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+1:], newOrderList.Orders[j:])
								newOrderList.Orders[j] = internalPressOrder
								return newOrderList
							} else if j == len(currentOrderList.Orders)-1 {
								fmt.Println("This third STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+2:], newOrderList.Orders[j+1:])
								newOrderList.Orders[j+1] = internalPressOrder
								return newOrderList
							}
						}
					}
				}
			}
		}
	} else {
		tempNum = definitions.N_FLOORS
		fmt.Println("You are going down")
		if currentOrderList.Orders[0].Floor == elevatorState.Destination { /* You can add in front of currentOrderList */
			fmt.Println("First order is the destination floor")
			newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
			copy(newOrderList.Orders[1:], newOrderList.Orders[:])
			newOrderList.Orders[0] = internalPressOrder
			return newOrderList
		} else { /* There are orders before destinationOrder */
			for i, order := range currentOrderList.Orders {
				fmt.Println("Order[", i, "]: ", order)
				if order.Floor < tempNum { // To check where you turn
					fmt.Println(" order.Floor > tempNum ")
					if order.Floor < internalPressOrder.Floor && elevatorState.LastFloor < internalPressOrder.Floor {
						fmt.Println("This IF STATEMENT")
						newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
						copy(newOrderList.Orders[i+1:], newOrderList.Orders[i:])
						newOrderList.Orders[i] = internalPressOrder
						return newOrderList
					}
					tempNum = order.Floor
				}
				if tempNum == elevatorState.Destination {
					fmt.Println("tempNum == elevatorState.Destination")

					for j, order2 := range currentOrderList.Orders {
						fmt.Println("Length, ", len(currentOrderList.Orders), ", j, ", j)
						if j > i {
							if order2.Floor > internalPressOrder.Floor {
								fmt.Println("The other IF STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+1:], newOrderList.Orders[j:])
								newOrderList.Orders[j] = internalPressOrder
								return newOrderList
							} else if j == len(currentOrderList.Orders)-1 {
								fmt.Println("This third STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+2:], newOrderList.Orders[j+1:])
								newOrderList.Orders[j+1] = internalPressOrder
								return newOrderList
							}
						}
					}
				}
			}
		}
	}
	return definitions.Orders{}
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

func KeepTrackOfAllAliveSlaves(updatedSlaveIdChan <-chan string, allSlavesMapChanMap map[string](chan map[string]bool)) {
	allSlavesMap := make(map[string]bool)
	deadTime := time.Second * 3 // 3 seconds and slave is assumed dead
	timerMap := make(map[string]*time.Timer)

	for {
		select {
		case slave_id := <-updatedSlaveIdChan:
			if allSlavesMap[slave_id] == false { // If slave is new
				fmt.Println("Creating new timer for:", slave_id)
				timerMap[slave_id] = time.NewTimer(deadTime)
			}
			fmt.Println("Slave ", slave_id, "is alive!")

			// We have to use a mutex, as maps are passed by reference
			// mutex.Lock()
			allSlavesMap[slave_id] = true
			// mutex.Unlock()

			fmt.Println("Resetting timer for:", slave_id)
			// Resetting timer, as slave is alive
			timerMap[slave_id].Stop()
			timerMap[slave_id].Reset(deadTime)
		default:
			for id, timer := range timerMap {
				select {
				case <-timer.C: // deadTime has passed
					fmt.Println("\nSlave ", id, "is assumed dead\n")
					// mutex.Lock()
					allSlavesMap[id] = false // Slave is assumed dead
					// mutex.Unlock()

				default: // Needed to avoid blocking of channels
					time.Sleep(time.Millisecond * 10)
				}
			}
		}
		// Send map of all slaves to all channels
		for key := range allSlavesMapChanMap {
			allSlavesMapChanMap[key] <- allSlavesMap
		}
	}
}

// func ElevatorAlive(aliveMessageFromElevatorFunctionsChanMap map[string]chan bool) {
// 	deadTime := time.Second * 3          // Time until reboot is needed
// 	killTimer := time.NewTimer(deadTime) // Initialize timer

// 	select {}
// 	for {
// 		if isAlive {
// 			killTimer.Stop()
// 			killTimer.Reset(deadTime)
// 		}
// 	}
// }
