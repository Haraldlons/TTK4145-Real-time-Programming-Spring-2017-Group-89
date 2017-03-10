package master

import (
	"../network"
	"definitions"
	"math"
	//"network"
)




func Run() bool {

	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newSlave.Run()
	check(err)

	totalOrderList := storage.GetOrderListFromFile()


	go network.ListenForUpdatesFromSlaves()
	go network.KeepTrackOfAliveSlaves() 	


		
	return true
}

func check(err error) {
	if err != nil {
		panic(err)
	}
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
