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
	"log"
	// "os"
	"time"
	// "fmt"
	"os/exec"
)

type Order struct {
	Floor     int
	Direction int
}

type Orders struct {
	Name []Order
}

func main() {
	fmt.Println("Main function started")
	time.Sleep(2 * time.Second)

	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newSlave.Run()
	if err != nil {
		log.Fatal(err)
	}
	// check(err)

	// interFace := make([]interface{}, 4)
	listOfNumbers := []int{1, 2, 3, 4, 43}
	secondListOfNumbers := []int{1, 2, 3, 4, 7}
	// fmt.Println("listOfNumbers", listOfNumbers)
	// fmt.Println("listOfNumbers[2]", listOfNumbers[2])
	// for i, s:= range listOfNumbers {
	// 	interFace[i] = Order{Floor: s, Direction: 1}
	// }
	totalOrderList := Orders{[]Order{{Floor: 2, Direction: -1}, {Floor: 3, Direction: 1}}}

	for i := range listOfNumbers {
		totalOrderList = Orders{append(totalOrderList.Name, Order{Floor: listOfNumbers[i], Direction: secondListOfNumbers[i]})}
	}

	// fmt.Println("interFace:", interFace)

	printEachOrder(totalOrderList, 2)
	fmt.Println("totalOrderList: ", totalOrderList)
	fmt.Println("Length of totalOrderList:", len(totalOrderList.Name))
	// SaveOrdersToFile(2, )

	for {
		time.Sleep(time.Second * 2)
	}

	return
}

func printEachOrder(orders Orders, length int) {
	// fmt.Println("Orders: ", orders)
	for i := 0; i < length; i++ {
		// fmt.Println("Hu")
		fmt.Println("i: ", i, "orders.Name[", i, "].Floor=", orders.Name[i].Floor)
		fmt.Println("i: ", i, "orders.Name[", i, "].Direction=", orders.Name[i].Direction)
	}
	fmt.Println("Done")
}

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
