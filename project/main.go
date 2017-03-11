package main

import (
	// "./src/definitions"
	//"./src/driver"
	// "./src/elevator"
	//"./src/network"
	//"./src/buttons"
	//"./src/driver"
	// "./src/storage"
	"./src/master"
	//"./src/watchdog"
	"fmt"
	// "log"
	// "os"
	//"time"
	// "os/exec"
)

func main() {
	fmt.Println("Main function started")
	master.Run()

} //End main
