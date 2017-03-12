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
	masterIsAliveChan := make(chan string, 1)

	stopListening := make(chan bool)

	go network.ListenAfterAliveMasterRegularly(masterIsAliveChan, stopListening)

	for {
		select {
		case master_id := <-masterIsAliveChan:
			fmt.Println("Master is still alive: , ", master_id)
		case <-time.After(time.Millisecond * 3000):
			fmt.Println("Master is not alive for the last three seconds")
			stopListening <- true
			fmt.Println("Has send stopListening signal to network.ListenAfterAliveMasterRegularly")
			masterHasDiedChan <- true
			return
		}
	}
}

func TakeInUpdatesInOrderListAndSendUpdatesOnChannels(updatedOrderList <-chan definitions.Orders, orderListForExecuteOrders chan<- definitions.Orders, completedCurrentOrder <-chan bool, orderListToExternalPresses chan<- definitions.Orders, elevator_id string, elevatorState definitions.ElevatorState) {

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
				network.SendUpdatesToMaster(msg, elevatorState)
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
