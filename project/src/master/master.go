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
	//elevatorOrders := []Orders{}
	return true
}

// Returns int corresponding to elevator with lowest cost (0:N_ELEVS-1)
func findLowestCostElevator(elevatorStates [definitions.N_ELEVS]definitions.ElevatorState, externalButtonPress definitions.Order) int {
	minCost := 2 * definitions.N_FLOORS
	bestElevator := 0
	destinationFloor := externalButtonPress.Floor
	destinationDirection := externalButtonPress.Direction

	for i := 0; i < definitions.N_ELEVS; i++ {
		travelDirection := findTravelDirection(elevatorStates[i].LastFloor, destinationFloor)
		tempCost := int(math.Abs(destinationFloor - elevatorStates[i].LastFloor))

		if elevatorStates[i].Destination == definitions.IDLE {
			// Elevator is idle
			tempCost = tempCost - 1 // Prioritize idle elevators

		} else if travelDirection != elevatorStates[i].Direction || travelDirection != destinationDirection {
			// Elevator already passed the destination or is moving in the wrong direction wrt. dest. Direction and travelDirection
			tempCost = tempCost + int(math.Abs(elevatorStates[i].Destination-elevatorStates[i].LastFloor))
		}

		if tempCost < minCost {
			minCost = tempCost
			bestElevator = i
		}
	}
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
