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
