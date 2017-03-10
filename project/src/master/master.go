package master

import (
	"definitions"
	"math"
	//"network"
)

func Initialize() bool {
	return true
}

func Run() bool {
	elevatorOrders := []Orders{
		
	}
	return true
}

// Finds the elevator closest to the destination floor decided by order.
// elevatorStates is a list of the states of every elevator
func findClosestElevator(order definitions.Order, elevatorStates [definitions.N_ELEVS]definitions.ElevatorState, idle [definitions.N_ELEVS]bool) int {
	closest := 0
	shortestDistance := definitions.N_FLOORS //Maximum distance to initialize variable

	for i := 0; i < definitions.N_ELEVS; i++ {
		distance := math.Abs(elevatorStates[i].LastFloor - order.Floor)

		if elevatorStates[i].Direction == order.Direction { // Elevators moving in the correct direction are evaluated first
			if order.floor == elevatorStates[i].LastFloor&idle[i] { //If elevator is on correct floor
				return i
			} else {
				if diff < shortestDistance {
					closest = i
				}
			}
		}
		else { //Elevator is moving in the opposite direction

		}
	}
	return closest
}

// Returns int corresponding to elevator with lowest cost (0:N_ELEV-1)
func findLowestCostElevator(elevatorStates [definitions.N_ELEVS]definitions.ElevatorState, destinationFloor int) int{
	cost := definitions.N_FLOORS
	bestElevator := 0
	for i:= 0; i < definitions.N_ELEVS; i++ { 
		direction := findDirection(elevatorStates[i].LastFloor, destinationFloor)

		if elevatorStates[i].Destination == definitions.IDLE { // Elevator is idle
			tempCost := math.Abs(destinationFloor - elevatorStates[i].LastFloor)
			if tempCost <= cost { //prioritize idle elevators
				cost = tempCost
				bestElevator = i
			}

		} else if direction == elevatorStates[i].Direction { // Elevator is moving in the correct direction
			tempCost := math.Abs(destinationFloor - elevatorStates[i].LastFloor)
			if tempCost < cost {
				cost = tempCost
				bestElevator = i
			}

		} else { // Elevator already passed the destination, or is moving in the wrong direction
			tempCost := math.Abs(destinationFloor - elevatorStates[i].LastFloor) + math.Abs(elevatorStates[i].Destination-elevatorStates[i].LastFloor)
			if tempCost < cost {
				cost = tempCost
				bestElevator = i
			}

		}
	}
	return bestElevator
}

func findDirection(lastFloor int, destinationFloor int) int {
	if destinationFloor > lastFloor {
		return definitions.DIR_UP
	else if destinationFloor == lastFloor
		return definitions.DIR_STOP
	else
		return definitions.DIR_DOWN
	}
}

func UpdateOrders(orders interface {}) {
	
	orders.Orders[i] = 
}