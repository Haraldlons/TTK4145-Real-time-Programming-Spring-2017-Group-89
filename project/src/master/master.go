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
	"time"
)

func Run() {
	fmt.Println("I'm a MASTER!")

	// Channel definitions
	totalOrderListChan := make(chan definitions.Elevators) // Channel for passing totalOrderList
	updateInAllSlavesMap := make(chan map[string]bool)     // Channel for passing updates to map containing all slaves

	time.Sleep(time.Second)
	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newSlave.Run()
	if err != nil {
	}

	// Various declarations
	// aliveSlavesList := []int{1, 2, 3}
	allSlavesMap := make(map[string]bool) // "true" implies slave is alive

	go network.ListenAfterAliveSlavesRegularly(&aliveSlavesList)
	go keepTrackOfAllAliveSlaves(updateInAllSlavesMap)

	go network.SendMasterIsAliveRegularly()

	go sendToSlavesOnUpdate(totalOrderListChan)

	// Needs allSlavesMap to be updated with all currently alive slaves
	redistributeOrders(allSlavesMap, totalOrderListChan) // Should only be ran on start-up. Depends on "sendToSlavesOnUpdate()"

	// "handleUpdatesFromSlaves" cannot be started before redistributeOrders() has returned
	go handleUpdatesFromSlaves(totalOrderListChan)

	for {
		time.Sleep(time.Second)
	}

	// listOfAliveSlaves := network.GetSlavesAlive()
	// redistributeOrders(&listOfAliveSlaves)
	// network.broadcastOrderlist(totalOrderList)

	// go KeepTrackOfAliveSlaves(&listOfAliveSlaves)

}

// Update order list in "orders" object with the command defined by externalButtonPress
func updateOrders(orders *definitions.Orders, externalButtonPress definitions.Order, elevatorState definitions.ElevatorState) {
	if checkForDuplicateOrder(orders, externalButtonPress) {
		fmt.Println("This order is already in the queue!")
		return
	}

	// Check to see if order should be placed first based on current elevator state
	if elevatorState.Direction == externalButtonPress.Direction && floorIsInbetween(orders.Orders[0].Floor, externalButtonPress.Floor, elevatorState.LastFloor, elevatorState.Direction) {
		// Insert Order in first position
		fmt.Println("Inserting order in first postion")

		orders.Orders = append(orders.Orders, definitions.Order{})
		copy(orders.Orders[1:], orders.Orders[:])
		orders.Orders[0] = externalButtonPress
		return
	}

	for i := 1; i < len(orders.Orders); i++ {
		direction := orders.Orders[i].Direction
		if externalButtonPress.Direction == direction { // Elevator is moving in the right direction
			switch direction {
			case definitions.DIR_UP:
				if externalButtonPress.Floor < orders.Orders[i].Floor {
					// Insert Order in position (i)
					fmt.Println("Inserting order in postion", i)

					orders.Orders = append(orders.Orders, definitions.Order{})
					copy(orders.Orders[i+1:], orders.Orders[i:])
					orders.Orders[i] = externalButtonPress
					return
				}
			case definitions.DIR_DOWN:
				if externalButtonPress.Floor > orders.Orders[i].Floor {
					// Insert Order in position (i+1)
					fmt.Println("Inserting order in postion", i)

					orders.Orders = append(orders.Orders, definitions.Order{})
					copy(orders.Orders[i+1:], orders.Orders[i:])
					orders.Orders[i] = externalButtonPress
					return

				}
			default:
				fmt.Println("Something weird is up, buddy")
			}
		}
	}
	// Place order at back of orderList
	fmt.Println("Placing order at back of order list")
	orders.Orders = append(orders.Orders, externalButtonPress)
}

func checkForDuplicateOrder(orders *definitions.Orders, externalButtonPress definitions.Order) bool {
	for i := range orders.Orders {
		if orders.Orders[i] == externalButtonPress {
			return true
		}
	}
	return false
}

func floorIsInbetween(orderFloor int, buttonFloor int, elevatorFloor int, direction int) bool {
	switch direction {
	case definitions.DIR_UP:
		return buttonFloor > elevatorFloor &&
			buttonFloor < orderFloor
	case definitions.DIR_DOWN:
		return buttonFloor < elevatorFloor &&
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
func findLowestCostElevator(elevatorStates map[string]definitions.ElevatorState, externalButtonPress definitions.Order) string {
	minCost := 2 * definitions.N_FLOORS
	bestElevator := "localhost"
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
			bestElevator = id
		}
		fmt.Println("Cost of elevator", id, ":", tempCost)
	}
	fmt.Println("Minimum cost:", minCost)
	return bestElevator
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

func handleUpdatesFromSlaves(totalOrderListChan chan definitions.Elevators) {
	msgChan := make(chan definitions.MSG_to_master)
	totalOrderList := definitions.Elevators{}
	// Initialize maps
	totalOrderList.OrderMap = make(map[string]definitions.Orders)
	totalOrderList.ElevatorStateMap = make(map[string]definitions.ElevatorState)

	// Start goroutine to listen for updates from slaves
	go network.ListenToSlave(msgChan)
	for {
		select {
		case msg := <-msgChan: // New message received
			// Update totalOrderList with information from message
			totalOrderList.OrderMap[msg.Id] = msg.Orders
			totalOrderList.ElevatorStateMap[msg.Id] = msg.ElevatorState

			// Get map of states
			elevatorStateMap := totalOrderList.ElevatorStateMap

			// Find elevator best suited for taking the received orders, and add orders to corresponding order lists
			for i := range msg.ExternalButtonPresses {
				elevator_id := findLowestCostElevator(elevatorStateMap, msg.ExternalButtonPresses[i])
				orders := totalOrderList.OrderMap[elevator_id]
				updateOrders(&orders, msg.ExternalButtonPresses[i], elevatorStateMap[elevator_id])
			}

			// fmt.Println("Total order list: ", totalOrderList)

			// Send updates to channel
			totalOrderListChan <- totalOrderList
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// When totalorderlist is updated, send to all slaves
func sendToSlavesOnUpdate(totalOrderListChan <-chan definitions.Elevators) {
	fmt.Println("Starting sending orders to slave")

	for {
		select {
		case totalOrderList := <-totalOrderListChan:
			// fmt.Println("Length of totalOrderlist: ", len(totalOrderList.OrderMap))
			if len(totalOrderList.OrderMap) != 0 {
				msg := definitions.MSG_to_slave{Elevators: totalOrderList}
				fmt.Println("Message sent from Master to slave:", msg)
				network.SendToSlave(msg)
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// Function to be ran when program is booting.
// Used to redistribute active orders of elevators that have died
func redistributeOrders(allSlavesMap map[string]bool, totalOrderListChan chan<- definitions.Elevators) {
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
						elevator_id := findLowestCostElevator(totalOrderList.ElevatorStateMap, orders[i])
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
