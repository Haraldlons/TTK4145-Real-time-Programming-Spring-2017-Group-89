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

func TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList <-chan definitions.Orders, orderListForExecuteOrders chan<- definitions.Orders, completedCurrentOrder <-chan bool, orderListToExternalPresses chan<- definitions.Orders, elevator_id string, updateElevatorStateForUpdatesInOrderList <-chan definitions.ElevatorState) {
	elevatorState := definitions.ElevatorState{}

	go func() {
		for {
			select {
			case elevatorState = <-updateElevatorStateForUpdatesInOrderList:
			}
		}
	}()

	currentOrderList := definitions.Orders{}
	storage.LoadOrdersFromFile(1, &currentOrderList)
	fmt.Println("Loaded totalOrderlist from a file. Result: ", currentOrderList)
	orderListForExecuteOrders <- currentOrderList
	time.Sleep(50 * time.Millisecond)

	// time.Sleep(50 * time.Millisecond)
	go func() {
		for {
			select {
			case <-completedCurrentOrder:
				if len(currentOrderList.Orders) > 0 {
					currentOrderList = definitions.Orders{currentOrderList.Orders[1:]}
				}
				orderListForExecuteOrders <- currentOrderList
				orderListToExternalPresses <- currentOrderList
				storage.SaveOrdersToFile(1, currentOrderList)
				msg := definitions.MSG_to_master{Orders: currentOrderList, Id: elevator_id}
				// fmt.Println("msg_to_master: ", msg)
				network.SendUpdatesToMaster(msg, elevatorState, elevator_id)
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	for {
		select {
		case updatedOrderList := <-updatedOrderList:
			//if currentOrderList != updatedOrderList { /*If orderlist from master is identical to our copy*/
			currentOrderList = updatedOrderList
			orderListForExecuteOrders <- currentOrderList
			orderListToExternalPresses <- currentOrderList
			storage.SaveOrdersToFile(1, currentOrderList)
			// sendOrderListUpdateToMaster(currentOrderList)
			time.Sleep(50 * time.Millisecond)
			// }
		}
	}
}

func slavesAlive(updatedSlaveIdChanMap map[string](chan string), allSlavesMapChanMap map[string](chan map[string]bool)) {
	allSlavesMap := make(map[string]bool)
	deadTime := time.Second * 3 // 3 seconds and slave is assumed dead
	// timer := time.NewTimer(deadTime)
	timerMap := make(map[string]*time.Timer)
	timerMapChan := make(chan map[string]*time.Timer)

	for {
		select {
		case slave_id := <-updatedSlaveIdChanMap["toWatchDog"]:
			for id := range timerMap {
				if timerMap[id] != nil {
					timerMap[id] = time.NewTimer(deadTime)
					timerMapChan <- timerMap
				}
			}

			timerMapChan[id]

		default:
			for id, timer := range timerMap {

			}

		}
		allSlavesMap[slave_id] = true

		// Send map of all slaves to channel
		allSlavesMapChanMap["toKeepTrackOfAllAliveSlaves"] <- allSlavesMap
	}

}
