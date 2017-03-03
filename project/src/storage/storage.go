package storage

import (
	"bufio"
	"fmt"
	//"io/ioutil"
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

	defer f.close()
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
	w := bufio.NewWriter(f)

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
	r := bufio.NewReader(f)
	return true
}

func StoreOrders(elevatorNum int) bool {
	return true
}

func LoadOrders(elvatorNum int) bool {
	return true
}
