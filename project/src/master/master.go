package master

import (

	"../network"
	"../definitions"
	"../storage"
	"math"
)


func Run() {

	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newSlave.Run()
	check(err)

	// Initialize Elevators struct to keep track of elevator orders
	var totalOrderList definitions.Elevators
	Elevators.OrderMap = make(map[string] []Orders)

	// Load from storage if available
	storage.LoadOrdersFromFile(&totalOrderList)

	listOfAliveSlaves := network.getSlavesAlive()
	redistributeOrders(&listOfAliveSlaves)
	network.broadcastOrderlist(totalOrderList)

	go handleUpdatesFromSlaves(&totalOrderList)
	go KeepTrackOfAliveSlaves(&listOfAliveSlaves)

}


func Run() bool {
	//elevatorOrders := []Orders{}
	return true
}

// Returns int corresponding to elevator with lowest cost (0:N_ELEVS-1)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleUpdatesFromSlaves(totalOrderList []definitions.Orders){

	totalOrderList := totalOrderList

	go network.listenForUpdatesFromSlave(totalOrderList)

	go func(){
		for {
			select {
				case <-totalOrderList
				// Handle updated orderList from slaves
			}
		}	
	}()
}

func KeepTrackOfAliveSlaves(&listOfAliveSlaves){

	AliveMessageFromSlave := make(chan slave)

	go network.listenAfterSlaves(AliveMessageFromSlave)
	select {
		case AliveMessageFromSlave := <-AliveMessageFromSlave
			for slave := range(listOfAliveSlaves){
				if AliveMessageFromSlave == listOfAliveSlaves[i]{
					slave<- "slave number" + slave "is alive"
				}
				
			}
	}


	for slave := range(listOfAliveSlaves){
		go func(){
			select {
				case <-slave
					fmt.Println("Slave:", slave, "is alive")
				case time.After(5*time.Second)
					fmt.Println("Slave:", slave, "died!")
					listOfAliveSlaves = listOfAliveSlaves.slice(deadSlave)
					redistributeOrders()
			}
		}()
	}
}


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
				fmt.Println("Elevator ", i+1, " has extra cost")
			}
		}

		if tempCost < minCost {
			minCost = tempCost
			bestElevator = i
		}
		fmt.Println("Cost of elevator", i+1, ":", tempCost)
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

func UpdateOrders(orders interface{}, externalButtonPress definitions.Order) {
	for i := range orders.Orders {
		direction := orders.Orders[i].Direction
		if externalButtonPress.Direction == direction { // Elevator is moving in the right direction
			switch direction {
			case definitions.DIR_UP:
				if externalButtonPress.Floor < orders.Orders[i].Floor {
					// Insert Order in position (i)
					orders.Orders = append(Orders[:i], append([]T{externalButtonPress}, orders.Orders[i:]...)...)
					return
				} else if externalButtonPress.Floor == orders.Orders[i].Floor {
					fmt.Println("Duplicate order in UpdateOrders()")
					return
				}
			case definition.DIR_DOWN:
				if externalButtonPress.Floor > orders.Orders[i].Floor {
					// Insert Order in position (i+1)
					orders.Orders = append(Orders[:i+1], append([]T{externalButtonPress}, orders.Orders[i+1:]...)...)
					return
				} else if externalButtonPress.Floor == orders.Orders[i].Floor {
					fmt.Println("Duplicate order in UpdateOrders()")
					return
				}
			default:
				//No clue
			}
		}
	}
	// Place order at back of orderList
	orders.Orders = append(orders.Orders, externalButtonPress)
}
