package master

import (
	"../definitions"
	// "../network"
	// "../storage"
	"fmt"
	"math"
	"os/exec"
)

func Run() {

	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newSlave.Run()
	check(err)

	// // Initialize Elevators struct to keep track of elevator orders
	// var totalOrderList definitions.Elevators
	// Elevators.OrderMap = make(map[string] []Orders)

	// // Load from storage if available
	// storage.LoadOrdersFromFile(&totalOrderList)

	// listOfAliveSlaves := network.GetSlavesAlive()
	// redistributeOrders(&listOfAliveSlaves)
	// network.broadcastOrderlist(totalOrderList)

	// go handleUpdatesFromSlaves(&totalOrderList)
	// go KeepTrackOfAliveSlaves(&listOfAliveSlaves)

}

func TestRun() {

	buttonPress := definitions.Order{Floor: 2, Direction: definitions.DIR_UP}
	stateList := []definitions.ElevatorState{
		definitions.ElevatorState{LastFloor: 0, Direction: definitions.DIR_UP, Destination: 1},
		definitions.ElevatorState{LastFloor: 1, Direction: definitions.DIR_UP, Destination: 2},
		definitions.ElevatorState{LastFloor: 1, Direction: definitions.DIR_UP, Destination: 2},
	}

	for i := range stateList {
		stateList[i].LastFloor = i
	}

	fmt.Println("Order: ", buttonPress)
	fmt.Println("Statelist:", stateList)

	bestElevator := findLowestCostElevator(stateList, buttonPress)
	fmt.Println("Best elevator: Elevator number ", bestElevator)

	orderList := []definitions.Order{
		definitions.Order{Floor: 2, Direction: definitions.DIR_DOWN},
		definitions.Order{Floor: 3, Direction: definitions.DIR_DOWN},
		definitions.Order{Floor: 4, Direction: definitions.DIR_UP},
		definitions.Order{Floor: 1, Direction: definitions.DIR_DOWN},
	}

	orders := definitions.Orders{
		Orders: orderList,
	}

	fmt.Println("Orders before update:", orders)
	state := definitions.ElevatorState{LastFloor: 4, Direction: definitions.DIR_DOWN, Destination: 0}
	updateOrders(&orders, buttonPress, state)
	fmt.Println("Orders after update:", orders)
}

/*
func handleUpdatesFromSlaves(totalOrderList chan definitions.Elevators) {
	go network.listenForUpdatesFromSlave(totalOrderList)
	go func() {
		for {
			select {
			case <-totalOrderList:
				bestElevator := findLowestCostElevator(elevatorStates, externalButtonPress)
				updateOrders(&orders, externalButtonPress, elevatorState)

			}
		}
	}()
}
*/

/*
func PrepareMessageToSlaves(slaveNum int, allElevators *definitions.Elevators) []byte {
	message, err := json.Marshal(allElevators)
	fmt.Println("JSON in ByteArray:", message)

	jsonByteLength := len(message)
	firstByte := jsonByteLength / 255
	secondByte := jsonByteLength - firstByte*255

	fmt.Println("JSONByteArrayLength:", jsonByteLength)

	fmt.Println(byte(len(message)))

	b = append([]byte{byte(secondByte)}, message...)
	b = append([]byte{byte(firstByte)}, message...)
}
*/

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

// Returns int corresponding to elevator with lowest cost (0:N_ELEVS-1)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// func handleUpdatesFromSlaves(totalOrderList []definitions.Orders){

// 	totalOrderList := totalOrderList

// 	go network.ListenForUpdatesFromSlave(totalOrderList)

// 	go func(){
// 		for {
// 			select {
// 				case <-totalOrderList
// 				// Handle updated orderList from slaves
// 			}
// 		}
// 	}()
// }

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

// Finds the elevator closest to the destination floor decided by order.
// elevatorStates is a list of the states of every elevator
func findClosestElevator(order definitions.Order, elevatorStates [definitions.N_ELEVS]definitions.ElevatorState, idle [definitions.N_ELEVS]bool) int {
	closest := 0
	shortestDistance := definitions.N_FLOORS //Maximum distance to initialize variable

	for i := 0; i < definitions.N_ELEVS; i++ {
		travelDirection := findTravelDirection(elevatorStates[i].LastFloor, destinationFloor)
		tempCost := int(math.Abs(float64(destinationFloor - elevatorStates[i].LastFloor)))

		if elevatorStates[i].Destination == definitions.IDLE {
			// Elevator is idle
			tempCost = tempCost - 1 // Prioritize idle elevators
		} else if elevatorStates[i].Destination != destinationFloor {
			// No additional cost if elevator destination is the same as order destination
			if elevatorHasAdditionalCost(travelDirection, destinationFloor, destinationDirection, elevatorStates[i]) {
				costToDest := int(math.Abs(float64(elevatorStates[i].Destination - elevatorStates[i].LastFloor)))
				tempCost = costToDest + int(math.Abs(float64(destinationFloor-elevatorStates[i].Destination)))
				fmt.Println("Elevator ", i, " has extra cost")
			}
		}

		if tempCost < minCost {
			minCost = tempCost
			bestElevator = i
		}
		fmt.Println("Cost of elevator", i, ":", tempCost)
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
