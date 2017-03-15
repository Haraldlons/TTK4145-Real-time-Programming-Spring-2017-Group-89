package storage

import (
	"../definitions"
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
	"encoding/json"
	"os"
	"strconv"
)

// TODO REMOVE THE ONES NOT NEEDED

const (
	FILEPATH                         = "./src/storage/"
	FILENAME_INTERNAL_BUTTON_PRESSES = "internal_button_presses.txt"
	FILENAME_EXTERNAL_BUTTON_PRESSES = "external_button_presses.txt"
	FILENAME_ELEVATOR_ORDERS         = "elevatorOrders"
	FILENAME_ELEVATORS               = "elevators.txt"
)

func checkError(err error) { // TODO REMOVE
	if err != nil {
		fmt.Print("Error in storage: ")
		panic(err)
	}
}

func SaveOrdersToFile(elevatorNum int, orders interface{}) { // TODO SEE IF NEEDED
	fileName := FILENAME_ELEVATOR_ORDERS + strconv.Itoa(elevatorNum)
	outFile, err := os.Create(FILEPATH + fileName + ".txt")
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(orders)
	checkError(err)
}

//Takes pointer as input arg
func LoadOrdersFromFile(elevatorNum int, orders interface{}) { // TODO SEE IF NEEDED
	fileName := FILENAME_ELEVATOR_ORDERS + strconv.Itoa(elevatorNum)
	inFile, err := os.Open(FILEPATH + fileName + ".txt")
	defer inFile.Close()
	checkError(err)

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(orders)
	checkError(err)
}

func SaveButtonPresses(typeOfButton string, buttonPresses interface{}) { // TODO SEE IF NEEDED
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

func LoadButtonPresses(typeOfButton string, buttonPresses interface{}) { // TODO SEE IF NEEDED
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


func SaveJSONtoFile(state interface{}) { //TODO DELETE
	// fmt.Println("Saving JSON to file")
	fileName := "JSON.txt"
	outFile, err := os.Create(FILEPATH + fileName)
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(state)
	checkError(err)
}


func SaveElevatorsToFile(elevators definitions.Elevators) {
	fmt.Println("Saving elevators to file")
	fileName := FILENAME_ELEVATORS
	outFile, err := os.Create(FILEPATH + fileName)
	defer outFile.Close()
	checkError(err)

	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(elevators)
	checkError(err)
}

func LoadElevatorsFromFile(elevators *definitions.Elevators) {
	fmt.Println("Loading elevators from file")
	fileName := FILENAME_ELEVATORS
	inFile, err := os.Open(FILEPATH + fileName)
	defer inFile.Close()
	checkError(err)

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(elevators)
	checkError(err)
}
