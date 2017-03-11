package master

import (
	"../definitions"
	"math"
	//"network"
)

func Initialize() bool {
	return true
}

func Run() bool {
	buttonPress := definitions.Order{Floor: 3, Direction: definitions.DIR_DOWN}
	stateList := []definitions.ElevatorState{
		definitions.ElevatorState{LastFloor: 0, Direction: definitions.DIR_UP, Destination: 1},
		definitions.ElevatorState{LastFloor: 1, Direction: definitions.DIR_UP, Destination: 2},
		definitions.ElevatorState{LastFloor: 2, Direction: definitions.DIR_UP, Destination: 3},
		definitions.ElevatorState{LastFloor: 3, Direction: definitions.DIR_UP, Destination: 4},
		definitions.ElevatorState{LastFloor: 4, Direction: definitions.DIR_DOWN, Destination: 0},
		definitions.ElevatorState{LastFloor: 5, Direction: definitions.DIR_UP, Destination: 6},
		definitions.ElevatorState{LastFloor: 6, Direction: definitions.DIR_UP, Destination: 7},
		definitions.ElevatorState{LastFloor: 7, Direction: definitions.DIR_UP, Destination: 8},
	}

	for i := range stateList {
		stateList[i].LastFloor = i
	}

	fmt.Println("Order: ", buttonPress)
	fmt.Println("Statelist:", stateList)

	bestElevator := findLowestCostElevator(stateList, buttonPress)
	fmt.Println("Best elevator: Elevator number ", bestElevator+1)
}

// Returns int corresponding to elevator with lowest cost (0:N_ELEVS-1)
func findLowestCostElevator(elevatorStates []definitions.ElevatorState, externalButtonPress definitions.Order) int {
	minCost := 2 * definitions.N_FLOORS
	bestElevator := 0
	destinationFloor := externalButtonPress.Floor
	destinationDirection := externalButtonPress.Direction

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
