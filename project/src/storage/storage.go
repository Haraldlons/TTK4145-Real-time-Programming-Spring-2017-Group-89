package storage

import (
	"../definitions"
	// "../driver"
	// "./../controller"
	// "./src/network"
	// "../buttons"
	//"./src/driver"
	// "../storage"
	//"./src/master"
	//"./src/watchdog"
	"bufio"
	"fmt"
	// "io/ioutil"
	"log"
	"os"
)

const (
	FILEPATH                         = "./src/storage/"
	FILENAME_INTERNAL_BUTTON_PRESSES = "internal_button_presses"
	FILENAME_EXTERNAL_BUTTON_PRESSES = "external_button_presses"
	FILENAME_ELEVATOR_ORDERS         = "elevatorOrders"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}



// func GetOrdersFromFile(elevatorNum int) (orders [definitions.ELEVATOR_ORDER_SIZE]definitions.Order) {
// 	orders = [definitions.ELEVATOR_ORDER_SIZE]definition.Order{}
// 	fileName := FILENAME_ELEVATOR_ORDERS

// 	// Open file
// 	file, _ := os.Open(fileName)
// 	defer file.Close()

// 	// Initialize reader object
// 	reader := bufio.NewReader(file)
// 	scanner := bufio.NewScanner(reader)
// 	scanner.Split(bufio.ScanWords)

// 	elevatorCount := 0 // Counter to keep track of which elevator's orders are being read
// 	i := 0
// 	for scanner.Scan() { // Scan every line
// 		line := scanner.Text()
// 		if line == "*" {
// 			elevatorFile++
// 			break
// 		}

// 		floor, _ := strconv.Atoi(line) // Convert string to int
// 		orders[i].floor = floor
// 		i++
// 	}

// 	for i := 0; i < 10; i++ {
// 		fmt.Println(orders[i].floor)
// 	}

// 	fmt.Println(orders)
// 	return orders
// }

// Harald spagetti code
func testFileWriting() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile("output.txt", os.O_WRONLY, 0666)
	defer outputFile.Close()

	var a, b int
	var itemCount int
	itemCount, err = fmt.Fscanf(inputFile, "%d %d\n", &a, &b)
	// fmt.Println("err: ", err.Error())
	for itemCount > 0 && err == nil {
		fmt.Println("itemCount: ", itemCount, "a: ", a, "b: ", b)
		fmt.Fprintln(outputFile, "B value is: ", b, ", and A value is: ", a)
		// fmt.Fprintln(w, ...)
		itemCount, err = fmt.Fscanln(inputFile, &a, &b)
	}
}

func ReadElevatorStateFromFile(elevatorState *definitions.ElevatorState) {
	inputFile, err := os.Open("output.txt") // output.txt is in the project folder. This still works
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	var lastFloor, direction int
	// var firstValue int
	firstValue, err := fmt.Fscanf(inputFile, "%d %d", &lastFloor, &direction)
	fmt.Println("Reading ElevatorState from File, lastFloor: ", lastFloor, ", direction: ", direction)
	// fmt.Println(elevatorState)
	firstValue++
	elevatorState.LastFloor = lastFloor
	elevatorState.Direction = direction
|
}

func SaveElevatorStateToFile(lastFloor int, direction int) {
	outputFile, err := os.OpenFile("output.txt", os.O_WRONLY, 0666) //This file is in project folder
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	fmt.Fprintln(outputFile, lastFloor, direction)
	fmt.Println("saving ElevatorState to File, lastFloor: ", lastFloor, ", direction: ", direction)
}

func SaveOrderToFile(order int) {
	outputFile, err := os.OpenFile("lastOrder.txt", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	fmt.Fprintln(outputFile, order)
	fmt.Println("Saving order to file. Order: ", order)
}

func getOrderFromFile() int {
	inputFile, err := os.Open("lastOrder.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	var order int
	// var uselessVariable int
	fmt.Fscanf(inputFile, "%d", &order)
	fmt.Println(", order: ", order)

	return order
}


// NEW FUNCTIONS

func SaveOrdersToFile(elevatorNum int, orders interface{}) {
	fileName := FILENAME + strconv.Itoa(elevatorNum)
	outFile, err := os.Create(fileName)
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(orders)
	checkError(err)
}

//Takes pointer as input arg
func LoadOrdersFromFile(elevatorNum int, orders interface{}) {
	fileName := FILENAME + strconv.Itoa(elevatorNum)
	inFile, err := os.Open(fileName)
	defer inFile.Close()
	checkError(err)

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(orders)
	checkError(err)
}

func SaveButtonPresses(typeOfButton string, buttonPresses interface{}) {
	fileName := "blank"
	if typeOfButton == "internal" {
		fileName = "internal_button_presses"
	} else if typeOfButton == "external" {
		fileName = "external_button_presses"
	} else {
		fmt.Println("Not a valid button type.")
	}

	outFile, err := os.Create(fileName)
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(buttonPresses)
	checkError(err)
}

func LoadButtonPresses(typeOfButton string, buttonPresses interface{}) {
	fileName := "blank"
	if typeOfButton == "internal" {
		fileName = "internal_button_presses"
	} else if typeOfButton == "external" {
		fileName = "external_button_presses"
	} else {
		fmt.Println("Not a valid button type.")
	}

	inFile, err := os.Open(fileName)
	defer inFile.Close()
	checkError(err)

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(buttonPresses)
	checkError(err)
}

func SaveStateToFile(state interface{}) {
	fileName := "state"
	outFile, err := os.Create(fileName)
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(state)
	checkError(err)
}

func LoadStateFromFile(state interface{}) {
	fileName := "state"
	inFile, err := os.Open(fileName)
	defer inFile.Close()
	checkError(err)

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(state)
	checkError(err)
}
