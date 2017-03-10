package storage

import (
	// "../definitions"
	// "../driver"
	// "../slave"
	// "./src/network"
	// "../buttons"
	//"./src/driver"
	// "../storage"
	//"./src/master"
	//"./src/watchdog"
	// "bufio"
	"fmt"
	// "io/ioutil"
	// "log"
	"os"
	"strconv"
	"encoding/json"
)

const (
	FILEPATH                         = "./src/storage/"
	FILENAME_INTERNAL_BUTTON_PRESSES = "internal_button_presses"
	FILENAME_EXTERNAL_BUTTON_PRESSES = "external_button_presses"
	FILENAME_ELEVATOR_ORDERS         = "elevatorOrders"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func SaveOrdersToFile(elevatorNum int, orders interface{}) {
	fileName := FILENAME_ELEVATOR_ORDERS + strconv.Itoa(elevatorNum)
	outFile, err := os.Create(fileName)
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(orders)
	checkError(err)
}

//Takes pointer as input arg
func LoadOrdersFromFile(elevatorNum int, orders interface{}) {
	fileName := FILENAME_ELEVATOR_ORDERS + strconv.Itoa(elevatorNum)
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

func SaveElevatorStateToFile(state interface{}) {
	fileName := "state"
	outFile, err := os.Create(FILEPATH + fileName)
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(state)
	checkError(err)
}

func LoadElevatorStateFromFile(state interface{}) {
	fileName := "state"
	inFile, err := os.Open(FILEPATH + fileName)
	defer inFile.Close()
	checkError(err)

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(state)
	checkError(err)
}
