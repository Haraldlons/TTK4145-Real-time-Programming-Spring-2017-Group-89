package watchdog

import (
	"../definitions"
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
}

func TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList <-chan definitions.Orders, orderListForExecuteOrders chan<- definitions.Orders, completedCurrentOrder <-chan bool, orderListToExternalPresses chan<- definitions.Orders, elevator_id string, updateElevatorStateForUpdatesInOrderList <-chan definitions.ElevatorState, orderListChanForPrinting chan<- definitions.Orders, lastSentMsgToMasterChanForPrinting chan<- definitions.MSG_to_master) {
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

	for {
		select {
		case currentOrderList := <-updatedOrderList:
			// if currentOrderList != updatedOrderList { /*If orderlist from master is identical to our copy*/
			fmt.Println("40")

			orderListForExecuteOrders <- currentOrderList
			fmt.Println("41")
			orderListToExternalPresses <- currentOrderList
			fmt.Println("42")
			orderListChanForPrinting <- currentOrderList
			fmt.Println("43")
			storage.SaveOrdersToFile(1, currentOrderList)
			// sendOrderListUpdateToMaster(currentOrderList)
			fmt.Println("40")

			time.Sleep(50 * time.Millisecond)
		case <-completedCurrentOrder:
			if len(currentOrderList.Orders) > 0 {
				currentOrderList = definitions.Orders{currentOrderList.Orders[1:]}
			}
			orderListForExecuteOrders <- currentOrderList
			fmt.Println("44")
			orderListToExternalPresses <- currentOrderList
			fmt.Println("45")
			storage.SaveOrdersToFile(1, currentOrderList)
			orderListChanForPrinting <- currentOrderList
			fmt.Println("46")
			msg := definitions.MSG_to_master{Orders: currentOrderList, Id: elevator_id}
			// fmt.Println("msg_to_master: ", msg)
			fmt.Println("47")
			network.SendUpdatesToMaster(msg, elevatorState, elevator_id, lastSentMsgToMasterChanForPrinting)
			fmt.Println("48")
			time.Sleep(50 * time.Millisecond)
		case elevatorState = <-updateElevatorStateForUpdatesInOrderList:

			// }
		}
	}
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
