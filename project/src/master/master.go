package master

import (
	"../definitions"
	// "../driver"
	"../network"
	"../storage"
	"fmt"
	"math"
	"os/exec"
	// "string"
	// "math/rand"
	// "net"
	"sync"
	"time"
)

func Run() {
	fmt.Println("I'm a MASTER!")

	master_id, _ := network.GetLocalIP()
	mutex := &sync.Mutex{}

	// Channel definitions
	totalOrderListChan := make(chan definitions.Elevators) // Channel for passing totalOrderList
	// updateInAllSlavesMap := make(chan map[string]bool)     // Channel for passing updates to map containing all slaves
	// allSlavesMap := make(map[string]bool)                  // "true" implies slave is alive

	// time.Sleep(time.Second)
	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run -race main.go")
	newSlave.Run()

	// Various declarations
	// aliveSlavesList := []int{1, 2, 3}

	//go network.ListenAfterAliveSlavesRegularly(&aliveSlavesList)
	// go keepTrackOfAllAliveSlaves(updateInAllSlavesMap)
	go network.SendMasterIsAliveRegularly(master_id)

	go sendToSlavesOnUpdate(totalOrderListChan, mutex)

	// Needs allSlavesMap to be updated with all currently alive slaves
	// redistributeOrders(allSlavesMap, totalOrderListChan, master_id) // Should only be ran on start-up. Depends on "sendToSlavesOnUpdate()"

	// "handleUpdatesFromSlaves" cannot be started before redistributeOrders() has returned
	go handleUpdatesFromSlaves(totalOrderListChan, master_id, mutex)

	for {
		time.Sleep(time.Second)
	}
}

// Update order list in "orders" object with the command defined by externalButtonPress
func updateOrders(orders *definitions.Orders, externalButtonPress definitions.Order, elevatorState definitions.ElevatorState) {
	if CheckForDuplicateOrder(orders, externalButtonPress.Floor) {
		fmt.Println("This order is already in the queue!")
		return
	}

	// fmt.Println("Orders received by updateOrders():", orders)
	// fmt.Println("Elevatorstate received by updateOrders():", elevatorState)
	// fmt.Println("ExternalButtonPress received by updateOrders():", externalButtonPress)

	if len(orders.Orders) > 0 { // For safety
		// Check to see if order should be placed first based on current elevator state
		if elevatorState.Direction == externalButtonPress.Direction && FloorIsInbetween(orders.Orders[0].Floor, externalButtonPress.Floor, elevatorState.LastFloor, elevatorState.Direction) {
			// Insert Order in first position
			// fmt.Println("Inserting order in first postion")

			orders.Orders = append(orders.Orders, definitions.Order{})
			copy(orders.Orders[1:], orders.Orders[:])
			orders.Orders[0] = externalButtonPress
			// fmt.Println("Orders returned by updateOrders():", orders)
			return
		}

	}

	for i := 1; i < len(orders.Orders); i++ {
		direction := orders.Orders[i].Direction
		if externalButtonPress.Direction == direction { // Elevator is moving in the right direction
			switch direction {
			case definitions.DIR_UP:
				if externalButtonPress.Floor < orders.Orders[i].Floor {
					// Insert Order in position (i)
					// fmt.Println("Inserting order in postion", i)
					orders.Orders = append(orders.Orders, definitions.Order{})
					copy(orders.Orders[i+1:], orders.Orders[i:])
					orders.Orders[i] = externalButtonPress
					// fmt.Println("Orders returned by updateOrders():", orders)
					return
				}
			case definitions.DIR_DOWN:
				if externalButtonPress.Floor > orders.Orders[i].Floor {
					// Insert Order in position (i+1)
					// fmt.Println("Inserting order in postion", i)

					orders.Orders = append(orders.Orders, definitions.Order{})
					copy(orders.Orders[i+1:], orders.Orders[i:])
					orders.Orders[i] = externalButtonPress
					// fmt.Println("Orders returned by updateOrders():", orders)
					return

				}
			default:
				fmt.Println("Something weird is up, buddy")
			}
		}
	}
	// Place order at back of orderList
	// fmt.Println("Placing order at back of order list")
	orders.Orders = append(orders.Orders, externalButtonPress)
	// fmt.Println("Orders returned by updateOrders():", orders)
}

// Don't accept more orders to same floor. Assume every person gets on elevator.
func CheckForDuplicateOrder(orders *definitions.Orders, buttonPressedFloor int) bool {
	for i := range orders.Orders {
		if orders.Orders[i].Floor == buttonPressedFloor {
			return true
		}
	}
	return false
}

func FloorIsInbetween(orderFloor int, buttonFloor int, elevatorLastFloor int, elevatorDirection int) bool {
	switch elevatorDirection {
	case definitions.DIR_UP:
		return buttonFloor > elevatorLastFloor &&
			buttonFloor < orderFloor
	case definitions.DIR_DOWN:
		return buttonFloor < elevatorLastFloor &&
			buttonFloor > orderFloor
	default:
		fmt.Println("Something is wrong in floorIsBetween()")
		return false
	}
}

// func KeepTrackOfAliveSlaves(&listOfAliveSlaves){

// 	AliveMessageFromSlave := make(chan slave)

// 	go network.listenAfterSlaves(AliveMessageFromSlave)
// 	select {
// 		case AliveMessageFromSlave := <-AliveMessageFromSlave
// 			for slave := range(listOfAliveSlaves){
// 				if AliveMessageFromSlave == listOfAliveSlaves[i]{
// 					slave<- "slave number" + slave "is alive"
// 				}

// 			}
// 	}

// 	for slave := range(listOfAliveSlaves){
// 		go func(){
// 			select {
// 				case <-slave
// 					fmt.Println("Slave:", slave, "is alive")
// 				case time.After(5*time.Second)
// 					fmt.Println("Slave:", slave, "died!")
// 					listOfAliveSlaves = listOfAliveSlaves.slice(deadSlave)
// 					redistributeOrders()
// 			}
// 		}()
// 	}
// }

// Returns id of best suited elevator, assuming elevatorStates is a map with ids
func findLowestCostElevator(elevatorStates map[string]definitions.ElevatorState, externalButtonPress definitions.Order, master_id string) string {
	minCost := 2 * definitions.N_FLOORS
	bestElevator_id := master_id
	destinationFloor := externalButtonPress.Floor
	destinationDirection := externalButtonPress.Direction

	for id, elevatorState := range elevatorStates { // Loop through map
		travelDirection := findTravelDirection(elevatorState.LastFloor, destinationFloor)
		tempCost := int(math.Abs(float64(destinationFloor - elevatorState.LastFloor)))

		if elevatorState.Destination == definitions.IDLE {
			// Elevator is idle
			tempCost = tempCost - 1 // Prioritize idle elevators
		} else if elevatorState.Destination != destinationFloor {
			// No additional cost if elevator destination is the same as order destination
			if elevatorHasAdditionalCost(travelDirection, destinationFloor, destinationDirection, elevatorState) {
				costToDest := int(math.Abs(float64(elevatorState.Destination - elevatorState.LastFloor)))
				tempCost = costToDest + int(math.Abs(float64(destinationFloor-elevatorState.Destination)))
				fmt.Println("Elevator with identifier", id, " has extra cost")
			}
		}

		if tempCost < minCost {
			minCost = tempCost
			bestElevator_id = id
		}
		// fmt.Println("Cost of elevator", id, ":", tempCost)
	}
	// fmt.Println("Minimum cost:", minCost)
	return bestElevator_id
}

func findTravelDirection(startFloor int, destinationFloor int) int {
	if destinationFloor > startFloor {
		return definitions.DIR_UP
	} else if destinationFloor == startFloor {
		return definitions.DIR_STOP
	} else {
		return definitions.DIR_DOWN
	}
}

// Returns true if elevator passes destinationFloor on it's way to elevatorDestination
func elevatorPassesDestinationFloor(travelDirection int, destinationFloor int, elevatorDestination int) bool {
	return (travelDirection == definitions.DIR_UP && destinationFloor-elevatorDestination < 0) ||
		(travelDirection == definitions.DIR_DOWN && destinationFloor-elevatorDestination > 0)
}

// Returns true if elevator can not go straight to destinationFloor
func elevatorHasAdditionalCost(travelDirection int, destinationFloor int, destinationDirection int, elevState definitions.ElevatorState) bool {
	return (elevatorPassesDestinationFloor(travelDirection, destinationFloor, elevState.Destination) &&
		travelDirection != destinationDirection) || // Elevator is traveling in the opposite direction of Order
		travelDirection != elevState.Direction || // Elevator is moving in the opposite direction relative to destination
		destinationFloor == elevState.LastFloor // Elevator has probably passed destination
}

func handleUpdatesFromSlaves(totalOrderListChan chan definitions.Elevators, elevator_id string, mutex *sync.Mutex) {
	msgChan := make(chan definitions.MSG_to_master)
	// completedUpdateOfOrderList := make(chan bool)
	// go func() {
	// completedUpdateOfOrderList <- true
	// }()

	totalOrderList := definitions.Elevators{}
	// Initialize maps
	totalOrderList.OrderMap = make(map[string]definitions.Orders)
	totalOrderList.ElevatorStateMap = make(map[string]definitions.ElevatorState)

	// Start goroutine to listen for updates from slaves
	go network.ListenToSlave(msgChan)

	for {
		// 	select {
		// 	case <-completedUpdateOfOrderList:
		select {
		case msg := <-msgChan: // New message received
			// fmt.Println("List of orders from msgChan:", msg.Orders)
			// fmt.Println("List of orders from totalOrderList:", totalOrderList.OrderMap[msg.Id])
			// Update totalOrderList with information from message
			fmt.Println("---------------------------------")
			mutex.Lock()
			fmt.Println("Message received from slave:", msg)
			mutex.Unlock()

			mutex.Lock()
			totalOrderList.OrderMap[msg.Id] = msg.Orders
			totalOrderList.ElevatorStateMap[msg.Id] = msg.ElevatorState
			mutex.Unlock()

			// Get map of states
			elevatorStateMap := totalOrderList.ElevatorStateMap

			// Find elevator best suited for taking the received orders, and add orders to corresponding order lists
			for i := range msg.ExternalButtonPresses {
				mutex.Lock()
				bestElevator_id := findLowestCostElevator(elevatorStateMap, msg.ExternalButtonPresses[i], elevator_id)
				mutex.Unlock()

				fmt.Println("Best elevator:", bestElevator_id, ", for order", msg.ExternalButtonPresses[i])

				mutex.Lock()
				orders := totalOrderList.OrderMap[bestElevator_id]
				updateOrders(&orders, msg.ExternalButtonPresses[i], elevatorStateMap[bestElevator_id])
				totalOrderList.OrderMap[bestElevator_id] = orders
				mutex.Unlock()

			}

			// fmt.Println("Total order list: ", totalOrderList)

			// Send updates to channel
			totalOrderListChan <- totalOrderList
			// go func() {
			// 	completedUpdateOfOrderList <- true
			// }()
			time.Sleep(time.Millisecond * 100)
		// case slavesAliveMap = <- slavesAliveMapToHandleUpdatesFromSlavesChan /*To be implemented*/
		}
	}
}

// When totalorderlist is updated, send to all slaves
func sendToSlavesOnUpdate(totalOrderListChan <-chan definitions.Elevators, mutex *sync.Mutex) {
	// fmt.Println("Starting sending orders to slave")

	for {
		select {
		case totalOrderList := <-totalOrderListChan:
			// fmt.Println("Length of totalOrderlist: ", len(totalOrderList.OrderMap))
			if len(totalOrderList.OrderMap) != 0 {
				msg := definitions.MSG_to_slave{Elevators: totalOrderList}
				fmt.Println("Message sent to slave:", msg)
				network.SendToSlave(msg, mutex)
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}

// Function to be ran when program is booting.
// Used to redistribute active orders of elevators that have died
func redistributeOrders(allSlavesMap map[string]bool, totalOrderListChan chan<- definitions.Elevators, master_id string) {
	defer fmt.Println("Orders have been redistributed and sent to network")
	totalOrderList := definitions.Elevators{}
	storage.LoadElevatorsFromFile(&totalOrderList)

	// Loop through the id of every currently alive slave
	for id_slaves, isAlive := range allSlavesMap {
		// Loop through maps of every elevator loaded from storage
		for id := range totalOrderList.OrderMap {
			if id_slaves == id && !isAlive { // Dead elevator
				orders := totalOrderList.OrderMap[id].Orders
				// Loop through every order
				for i := range orders {
					if orders[i].Direction != 0 { // Not an internal order
						// Find elevator with lowest cost, and add order to corresponding orderList
						elevator_id := findLowestCostElevator(totalOrderList.ElevatorStateMap, orders[i], master_id)
						updatedOrders := totalOrderList.OrderMap[id]
						updateOrders(&updatedOrders, orders[i], totalOrderList.ElevatorStateMap[elevator_id])
						totalOrderList.OrderMap[id] = updatedOrders
					}
				}
			}
		}
	}

	// Send updates to channel, which in turn is sent over network
	totalOrderListChan <- totalOrderList
}
