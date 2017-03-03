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
	//"io/ioutil"
	"log"
	"os"
)

const (
	FILEPATH                         = "./src/storage/"
	FILENAME_INTERNAL_BUTTON_PRESSES = "internal_button_presses"
	FILENAME_EXTERNAL_BUTTON_PRESSES = "external_button_presses"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//Internal button presses:
func StoreInternalButtonPresses() bool {
	fileName := FILENAME_INTERNAL_BUTTON_PRESSES

	f, err := os.Create(FILEPATH + fileName)
	check(err) //Remove when testing is done
	if err != nil {
		return false
	}

	defer f.Close()
	w := bufio.NewWriter(f)

	buttonPresses := []byte{'A', 'B', 'C', 'D'}

	n1, err := w.WriteString("Internal button presses:\n")
	check(err)
	if err != nil {
		return false
	}

	fmt.Printf("Wrote %d bytes\n", n1)
	w.Write(buttonPresses)

	w.Flush()
	return true
}

func LoadInternalButtonPresses() bool {
	fileName := FILENAME_INTERNAL_BUTTON_PRESSES

	f, err := os.Create(FILEPATH + fileName)
	check(err) //remove later
	if err != nil {
		return false
	}

	defer f.Close()
	r1 := bufio.NewReader(f)

	b1, err := r1.Peek(5)
	check(err) //remove later
	if err != nil {
		return false
	}
	fmt.Printf("5 bytes: %s\n", string(b1))

	return true
}

//External button presses:
func StoreExternalButtonPresses() bool {
	fileName := FILENAME_EXTERNAL_BUTTON_PRESSES

	f, err := os.Create(FILEPATH + fileName)
	check(err) //remove later
	if err != nil {
		return false
	}
	defer f.Close()
	// w := bufio.NewWriter(f)

	return true
}

func LoadExternalButtonPresses() bool {
	fileName := FILENAME_EXTERNAL_BUTTON_PRESSES

	f, err := os.Open(FILEPATH + fileName)
	check(err) //remove later
	if err != nil {
		return false
	}
	defer f.Close()
	// r := bufio.NewReader(f)
	return true
}

func StoreOrders(elevatorNum int) bool {
	return true
}

func LoadOrders(elvatorNum int) bool {
	return true
}

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
