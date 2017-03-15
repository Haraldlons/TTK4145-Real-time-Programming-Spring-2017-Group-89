package watchdog

import (
	"../def"
	// "../master"
	"../network"
	"../storage"
	// "net"
	"fmt"
	"time"
	"os"
	"../driver"
	"os/exec"
)

// var timeLimit = 10 * time.Second
// var listenTimer = 100 * time.Millisecond

// func SendNetworkAlive(udpBroadcast *net.UDPConn) bool { // Possibly unnesces
// 	msg := make([]byte, 1024)
// 	numLines, err := udpBroadcast.Write(msg)
// 	return err != nil
// }

func CheckIfMasterIsAliveRegularly(/*stopListening map[string] chan bool*/) {
	masterIsAliveChan := make(chan string)
	stopListening := make(chan bool)

	go network.ListenAfterAliveMasterRegularly(masterIsAliveChan, stopListening)
	for {
		select {
		case /*master_id := */<-masterIsAliveChan:

		// fmt.Println("Master is still alive: , ", master_id)

		case <-time.After(time.Second * 3):
			fmt.Println("Master is not alive for the last three seconds")
			driver.Elev_set_motor_direction(def.DIR_STOP)

			// Kill all network processes
			// stopListening <- true
			time.Sleep(time.Second*2)

			// Spawn new master
			newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			newSlave.Run()
			time.Sleep(time.Second*2)
			os.Exit(1)
			return
		}
	}

	// if master_id == "" {
	// }
}

func TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList <-chan def.Orders, orderListForExecuteOrders chan<- def.Orders, completedCurrentOrder <-chan bool, elevator_id string, orderListChanForPrinting chan<- def.Orders, lastSentMsgToMasterChanForPrinting chan<- def.MSG_to_master, orderListForSendingToMaster chan def.Orders, sendMessageToMaster chan bool, newInternalButtonOrderChan chan def.Order) {

	currentOrderList := def.Orders{}
	storage.LoadOrdersFromFile(1, &currentOrderList)
	fmt.Println("Loaded totalOrderlist from a file. Result: ", currentOrderList)
	orderListForExecuteOrders <- currentOrderList
	// newInternalButtonPress := def.Order{}
	lastOrderList := def.Orders{}

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

		case <-completedCurrentOrder:
			fmt.Printf("completedCurrentOrder")
			fmt.Println("CurrentOrderlist in special case:", currentOrderList)
			if len(currentOrderList.Orders) > 0 {
				fmt.Println("completedCurrentOrder23")
				currentOrderList = def.Orders{currentOrderList.Orders[1:]}
				fmt.Println("orderListafterSlice: ", currentOrderList)
			}
			orderListForExecuteOrders <- currentOrderList
			fmt.Println("45")
			fmt.Println("46")
			storage.SaveOrdersToFile(1, currentOrderList)
			fmt.Println("46.5")
			orderListChanForPrinting <- currentOrderList
			fmt.Println("47")
			// msg := def.MSG_to_master{Orders: currentOrderList, Id: elevator_id}
			// fmt.Println("msg_to_master: ", msg)
			fmt.Println("48")
			// network.SendUpdatesToMaster(msg, lastSentMsgToMasterChanForPrinting)
			fmt.Println("49")
			orderListForSendingToMaster <- currentOrderList
			sendMessageToMaster <- true

			fmt.Println("50")
		// case newInternalButtonPress = <-newInternalButtonOrderChan:
		// 	fmt.Println("50,5")
		// 	currentOrderList = distributeInternalOrderToOrderList(newInternalButtonPress, currentOrderList, elevatorState)
		// 	fmt.Println("51")
		// 	orderListForExecuteOrders <- currentOrderList
		// 	fmt.Println("52")
		// 	orderListChanForPrinting <- currentOrderList
		// 	fmt.Println("53")
		// 	orderListForSendingToMaster <- currentOrderList
		// 	fmt.Println("54")
		}
	}
}


func checkIfChangedOrderList(lastOrderList def.Orders, currentOrderList def.Orders) bool {
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


func KeepTrackOfAllAliveSlaves(updatedSlaveIdChan <-chan string, allSlavesAliveMapChanMap map[string](chan map[string]bool)) {
	allSlavesAliveMap := make(map[string]bool)
	deadTime := time.Second * 3 // 3 seconds and slave is assumed dead
	timerMap := make(map[string]*time.Timer)

	for {
		select {
		case slave_id := <-updatedSlaveIdChan:
			if allSlavesAliveMap[slave_id] == false { // If slave is new
				fmt.Println("Creating new timer for:", slave_id)
				timerMap[slave_id] = time.NewTimer(deadTime)
			}
			// fmt.Println("Slave ", slave_id, "is alive!")

			// We have to use a mutex, as maps are passed by reference
			// mutex.Lock()
			allSlavesAliveMap[slave_id] = true
			// mutex.Unlock()

			// fmt.Println("Resetting timer for:", slave_id)
			// Resetting timer, as slave is alive
			timerMap[slave_id].Stop()
			timerMap[slave_id].Reset(deadTime)
		default:
			for id, timer := range timerMap {
				select {
				case <-timer.C: // deadTime has passed
					// fmt.Println("\nSlave ", id, "is assumed dead\n")
					// mutex.Lock()
					allSlavesAliveMap[id] = false // Slave is assumed dead
					// mutex.Unlock()

				default: // Needed to avoid blocking of channels
					time.Sleep(time.Millisecond * 10)
				}
			}
		}
		// Send map of all slaves to all channels
		for key := range allSlavesAliveMapChanMap {
			// fmt.Println("allSlavesAliveMap in watchdog:", allSlavesAliveMap)
			allSlavesAliveMapChanMap[key] <- allSlavesAliveMap
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
