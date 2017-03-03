package main

import (
	//"master"
	"./src/definitions"
	"./src/driver"
	// "./src/network"
	//"./src/buttons"
	//"./src/driver"
	// "./src/storage"
	//"./src/master"
	//"./src/watchdog"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
	// "fmt"
	"net"
	// "os/exec"
)

var delay = 50 * time.Millisecond
var endProgram = false
var elevatorActive = false
var elevatorState = definitions.ElevatorState{2, 0}
var goToFloorVariable = -1
var msg = make([]byte, 8)

func main() {
	fmt.Println("Main function started")
	//network.Run()

	driver.Elev_init()
	driver.Elev_set_motor_direction(driver.DIRECTION_STOP)

	// elevatorState := definitions.ElevatorState{2, 0}
	readElevatorStateFromFile(&elevatorState)
	fmt.Println("elevatorInfo during initialization: ", elevatorState)

	stopSignal := 0
	// buttonSignal := driver.Elev_get_button_signal(0,0)
	if goToFloorVariable != -1 {
		go goToFloor(goToFloorVariable, &elevatorState)
	}

	driver.Elev_set_floor_indicator(3)
	goToFirstFloor := 0
	goToSecondFloor := 0
	goToThirdFloor := 0
	goToFourthFloor := 0

	// go goToFloor(3, &elevatorState)

	// retrieveElevatorStateFromFile()
	// testFileWriting()
	if getOrderFromFile() != -1 {
		go goToFloor(getOrderFromFile(), &elevatorState)
	}

	go setupNetwork()

	for {
		// fmt.Println("Elev_get_floor_sensor_signal: ", driver.Elev_get_floor_sensor_signal())
		printLastFloorIfChanged(&elevatorState)
		// updateElevatorStateIfChanged(&elevatorState)

		//driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)

		// if driver.Elev_get_floor_sensor_signal() == driver.N_FLOORS - 1 {
		// 	// fmt.Println("Bobby Brown")
		// 	driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
		// } else if driver.Elev_get_floor_sensor_signal() == 0{
		// 	// fmt.Println("Bobby Brown inverse")
		// 	driver.Elev_set_motor_direction(driver.DIRECTION_UP)
		// }

		goToFirstFloor = driver.Elev_get_button_signal(2, 0)
		goToSecondFloor = driver.Elev_get_button_signal(2, 1)
		goToThirdFloor = driver.Elev_get_button_signal(2, 2)
		goToFourthFloor = driver.Elev_get_button_signal(2, 3)

		if goToFirstFloor == 1 {
			// go goToFloor(0, &elevatorState)
			setOrderOverNetwork(0)
		}
		if goToSecondFloor == 1 {
			// go goToFloor(1, &elevatorState)
			setOrderOverNetwork(1)
		}
		if goToThirdFloor == 1 {
			// go goToFloor(2, &elevatorState)
			setOrderOverNetwork(2)
		}
		if goToFourthFloor == 1 {
			// go goToFloor(3, &elevatorState)
			setOrderOverNetwork(3)
		}

		if endProgram {
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			fmt.Println("endProgram == true. Stopping program")
			return
		}

		stopSignal = driver.Elev_get_stop_signal()
		if stopSignal != 0 {
			setOrderOverNetwork(0)
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			fmt.Println("Stopping program, with stop signal: ", stopSignal)
			fmt.Println("Another call to Elev_get_stop_signal(): ", driver.Elev_get_stop_signal())
			return
		}
	}
}

/*This functions should be cleaned up. I have an ide how to do it*/
func printLastFloorIfChanged(elevatorState *definitions.ElevatorState) {
	lastFloor := driver.Elev_get_floor_sensor_signal()
	switch lastFloor {
	case 0:
		if elevatorState.LastFloor != 0 {
			elevatorState.Direction = definitions.DIR_UP
			elevatorState.LastFloor = 0
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
			fmt.Println("Last Floor: 1. Direction: ", elevatorState.Direction)
		}
	case 1:
		if elevatorState.LastFloor != 1 {
			if elevatorState.LastFloor > lastFloor {
				elevatorState.Direction = definitions.DIR_DOWN
			} else {
				elevatorState.Direction = definitions.DIR_UP
			}
			elevatorState.LastFloor = 1
			fmt.Println("Last Floor: 2. Direction: ", elevatorState.Direction)
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
		}
	case 2:
		if elevatorState.LastFloor != 2 {
			if elevatorState.LastFloor > lastFloor {
				elevatorState.Direction = definitions.DIR_DOWN
			} else {
				elevatorState.Direction = definitions.DIR_UP
			}

			elevatorState.LastFloor = 2
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
			fmt.Println("Last Floor: 3. Direction: ", elevatorState.Direction)
		}
	case 3:
		if elevatorState.LastFloor != 3 {
			elevatorState.Direction = definitions.DIR_DOWN
			elevatorState.LastFloor = 3
			fmt.Println("Last Floor: 4. Direction: ", elevatorState.Direction)
			saveElevatorStateToFile(elevatorState.LastFloor, elevatorState.Direction)
		}

	default:

	}
}

func goToFloor(destinationFloor int, elevatorState *definitions.ElevatorState) {
	if !elevatorActive {
		saveOrderToFile(destinationFloor)
		elevatorActive = true

		fmt.Println("Going to floor: ", destinationFloor+1)
		direction := elevatorState.Direction
		lastFloor := elevatorState.LastFloor

		if driver.Elev_get_floor_sensor_signal() == destinationFloor {
			saveOrderToFile(-1)

			fmt.Println("You are allready on the desired floor")
			elevatorActive = false
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			// endProgram = true
			return
		} else { /*You are not on the desired floor*/
			driver.Elev_set_door_open_lamp(0)
			if lastFloor == destinationFloor {
				if direction == 1 {
					driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
				} else {
					driver.Elev_set_motor_direction(driver.DIRECTION_UP)
				}
			}
			if lastFloor < destinationFloor {
				driver.Elev_set_motor_direction(driver.DIRECTION_UP)
			} else {
				driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
			}
			for {
				if driver.Elev_get_floor_sensor_signal() == destinationFloor {
					time.Sleep(time.Millisecond * 150) //So the elevator stops in the middle of the sensor
					fmt.Println("You reached your desired floor. Walk out\n")
					elevatorActive = false
					// driver.Elev_set_button_lamp(1,1,1)
					// driver.Elev_set_button_lamp(0,1,1)
					driver.Elev_set_floor_indicator(destinationFloor)
					driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
					// endProgram = true
					time.Sleep(delay * 10)
					driver.Elev_set_door_open_lamp(1)
					saveOrderToFile(-1)
					return
				} else if driver.Elev_get_floor_sensor_signal() == 0 { /*This is just to be fail safe*/
					driver.Elev_set_motor_direction(driver.DIRECTION_UP)
				} else if driver.Elev_get_floor_sensor_signal() == 3 {
					driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
				} else {
					time.Sleep(delay)
				}
			}
		}
	}
}

func setFloorIndicator() {
	sensorValue := driver.Elev_get_floor_sensor_signal()
	if sensorValue != -1 {
		driver.Elev_set_floor_indicator(sensorValue)
	}
}

// func writeElevatorStateToFile(elevatorState definitions.elevatorState) {

// }

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

func readElevatorStateFromFile(elevatorState *definitions.ElevatorState) {
	inputFile, err := os.Open("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	var lastFloor, direction int
	var firstValue int
	firstValue, err = fmt.Fscanf(inputFile, "%d %d", &lastFloor, &direction)
	fmt.Println("firstValue: ", firstValue, ", lastFloor: ", lastFloor, ", direction: ", direction)
	fmt.Println(elevatorState)
	elevatorState.LastFloor = lastFloor
	elevatorState.Direction = direction

}

func saveElevatorStateToFile(lastFloor int, direction int) {
	outputFile, err := os.OpenFile("output.txt", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	fmt.Fprintln(outputFile, lastFloor, direction)
	fmt.Println(outputFile, lastFloor, direction)
}

func saveOrderToFile(order int) {
	outputFile, err := os.OpenFile("lastOrder.txt", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	fmt.Fprintln(outputFile, order)
	fmt.Println(outputFile, order)
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

/*Below is experimental work done a friday afternoon*/

// import (
// 	"encoding/binary"
// 	// "fmt"
// 	"net"
// 	"os/exec"
// 	// "time"
// )

var bcAddress string = "129.241.187.255"
var port string = ":55555"

// var delay = 300 * time.Millisecond

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func slave(udpListen *net.UDPConn) int {
	listenChan := make(chan int, 1)
	slaveCount := 0

	// Run goroutine listening for sent values from master.
	go listen(listenChan, udpListen)

	for {
		select {
		case slaveCount = <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("Got listen message: ", slaveCount)
			if slaveCount < 4 && slaveCount > -1 {
				fmt.Println("Going to floor from slave: ", slaveCount)
				go goToFloor(slaveCount, &elevatorState)
			}
			time.Sleep(delay / 2) // wait 50 ms
			break
		case <-time.After(100 * delay): // Wait 10 cycles (1 second). Master assumed dead
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			fmt.Println("Master is dead. Long live the the new king!")
			return slaveCount
		}
	}
}

func master(startCount int, udpBroadcast *net.UDPConn) {
	/* Launch new instance of "main".
	 * This creates the corresponding slave which will loop on listen until master dies
	 */
	// newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	// err := newSlave.Run()
	// check(err)

	count := startCount

	for {
		// Convert count from int to binary/byte and place in msg
		binary.BigEndian.PutUint64(msg, uint64(count))
		udpBroadcast.Write(msg)

		// fmt.Println(count)
		// count++

		time.Sleep(10 * delay) // Wait 1 cycle (100 ms)
	}
}

func listen(listenChan chan int, udpListen *net.UDPConn) {
	buf := make([]byte, 8)
	for {
		udpListen.ReadFromUDP(buf)

		// Convert byte from buf to int and send over channel.
		listenChan <- int(binary.BigEndian.Uint64(buf))
		time.Sleep(delay) // Wait 1 cycle (100 ms)
	}
}

func setupNetwork() {

	udpAddr, err := net.ResolveUDPAddr("udp", port)
	check(err)

	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)

	// Initialize slave
	// First run of program will return 0 and initialize master->slave topology
	fmt.Println("Run slave")
	count := slave(udpListen)

	udpListen.Close()

	udpAddr, err = net.ResolveUDPAddr("udp", bcAddress+port)
	check(err)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	count = 10
	fmt.Println("Run master")
	master(count, udpBroadcast)

	fmt.Println("Close broadcast")
	udpBroadcast.Close()
}

func setOrderOverNetwork(destinationFloor int) {
	fmt.Println("Sending order over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+port)
	check(err)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	binary.BigEndian.PutUint64(msg, uint64(destinationFloor))
	udpBroadcast.Write(msg)
}
