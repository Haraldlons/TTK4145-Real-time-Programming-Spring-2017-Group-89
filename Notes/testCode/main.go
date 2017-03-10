package main

import (
	// "./src/controller"
	// "./src/definitions"
	// "./src/driver"
	// "./src/elevator"
	// "./src/network"
	//"./src/buttons"
	//"./src/driver"
	// "./src/storage"
	//"./src/master"
	//"./src/watchdog"
	"fmt"
	// "log"
	// "os"
	// "time"
	// "fmt"
	// "os/exec"
)

type Order struct {
	Floor     int
	Direction int
}

type Orders struct {
	Orders []Order
}


func main() {
	fmt.Println("Main function started")

	interFace := Orders{
		[
		Order{Floor: 2, Direction: -1}, 
		Order{Floor: 3, Direction: 1},
	],
	}
	fmt.Println("interface: ", interFace)

	// SaveOrdersToFile(2, )

	return
} //End main

// func SaveOrdersToFile(elevatorNum int, orders interface{}) {
// 	fileName := FILENAME_ELEVATOR_ORDERS + strconv.Itoa(elevatorNum)
// 	outFile, err := os.Create(fileName)
// 	defer outFile.Close()
// 	checkError(err)

// 	encoder := json.NewEncoder(outFile)
// 	err = encoder.Encode(orders)
// 	checkError(err)
// }

// func LoadOrdersFromFile(elevatorNum int, orders interface{}) {
// 	fileName := FILENAME_ELEVATOR_ORDERS + strconv.Itoa(elevatorNum)
// 	inFile, err := os.Open(fileName)
// 	defer inFile.Close()
// 	checkError(err)

// 	decoder := json.NewDecoder(inFile)
// 	err = decoder.Decode(orders)
// 	checkError(err)
// }

// type ElevatorState struct {
// 	LastFloor int
// 	Direction int
// 	Destination int 
// }


// type Elevators struct {
// 	OrderMap map[string][]Orders
// }


// type MSG_to_master struct {
// 	Orders Orders
// 	ElevatorState ElevatorState
//  	ExternalButtonPresses []Order
//  	Id string
// }

// type MSG_to_slave struct {
// 	Elevators Elevators
// }




// type TestMessage struct {
//     Name string
//     Body string
//     Time int64
// }