package network

import (
	"../definitions"
	// "../driver"
	// "../buttons"
	"../storage"
	// "../elevator"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
	// "math"
	//"bufio"
	// "log"
	//"os"
	//"strconv"
)

var bcAddress string = "129.241.187.151"

//var bcAddress string = "localhost"
var port string = ":46723"
var slaveIsAlivePort string = ":46720"
var masterIsAlivePort string = ":46721"
var jsonSendPort string = ":46722"

var delay100ms = 100 * time.Millisecond

func check(err error) {
	if err != nil {
		fmt.Print("Error in network: ")
		panic(err)
	}
}

func slave(udpListen *net.UDPConn) int {
	defer fmt.Println("Slave function ended") /*For debugging*/
	listenChan := make(chan int, 1)
	slaveCount := 0

	stopSlaveBroadcasting := make(chan int)

	// Run goroutine listening for sent values from master.
	go listen(listenChan, udpListen)
	// go sendImAliveMessage(stopSlaveBroadcasting)

	for {
		select {
		case slaveCount = <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("Got listen message from master: ", slaveCount)
			// if slaveCount < 4 && slaveCount > -1 {
			// 	fmt.Println("Going to floor from slave: ", slaveCount)
			// 	// go goToFloor(slaveCount, &elevatorState)
			// }
			time.Sleep(delay100ms / 2) // wait 50 ms
			break
		case <-time.After(30 * delay100ms): // Wait 10 cycles (1 second). Master assumed dead
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			fmt.Println("Master is dead. Long live the the new king!")
			stopSlaveBroadcasting <- 1
			return slaveCount
		}
	}

}

func master(startCount int, udpBroadcast *net.UDPConn) {
	/* Launch new instance of "main".
	 * This creates the corresponding slave which will loop on listen until master dies
	 */
	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newSlave.Run()
	check(err)
	msg := make([]byte, 8)

	count := startCount

	// go listenAfterImAliveMessage()

	for {
		// Convert count from int to binary/byte and place in msg
		binary.BigEndian.PutUint64(msg, uint64(count))
		udpBroadcast.Write(msg)

		// fmt.Println(count)
		count++

		time.Sleep(10 * delay100ms) // Wait 1 cycle (100 ms)
	}
}

func listen(listenChan chan int, udpListen *net.UDPConn) {
	buf := make([]byte, 8)
	for {
		udpListen.ReadFromUDP(buf)

		// Convert byte from buf to int and send over channel.
		listenChan <- int(binary.BigEndian.Uint64(buf))
		time.Sleep(10 * delay100ms) // Wait 1 cycle (100 ms)
	}
}

func SetupNetwork() {

	// sendJSON()

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
	// go RecieveJSON()
	master(count, udpBroadcast)

	fmt.Println("Close broadcast")
	udpBroadcast.Close()
}

func SendSlaveIsAliveRegularly(slaveID int, stopSlaveBroadcasting chan bool) {
	fmt.Println("Sending ImAliveMessage over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+slaveIsAlivePort)
	check(err)
	msg := make([]byte, 8)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	defer func() {
		fmt.Println("Actually stopping sending sendImAliveMessage")
		udpBroadcast.Close()
	}()

	binary.BigEndian.PutUint64(msg, uint64(slaveID))
	// go func() {
	// 	select {
	// 	case <- stopSlaveBroadcasting:
	// 		fmt.Println("Ending sendImAliveMessage")
	// 		return
	// 	}
	// }()
	for {
		select {
		case <-stopSlaveBroadcasting:
			time.Sleep(10 * time.Millisecond)
			return
		default:
			udpBroadcast.Write(msg)
			// fmt.Println("Sending I'm Alive")
			time.Sleep(delay100ms * 20)
		}
	}
}

func ListenAfterAliveSlavesRegularly(aliveSlavesList *[]int) {
	fmt.Println("Listening after ImAlive messages, list: ", aliveSlavesList)
	udpAddr, err := net.ResolveUDPAddr("udp", slaveIsAlivePort)
	check(err)

	slaveMessagesRecieved := 0

	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)
	defer udpListen.Close()

	listenChan := make(chan int, 1)
	// slaveCount := 0

	go func() {
		buf := make([]byte, 8)
		for {
			udpListen.ReadFromUDP(buf)

			// Convert byte from buf to int and send over channel.
			listenChan <- int(binary.BigEndian.Uint64(buf))
			time.Sleep(delay100ms) // Wait 1 cycle (100 ms)
		}
	}()
	notRecievedCounter := 0
	for {
		select {
		case <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			// fmt.Println("ListenAfterSlaves: ", slaveCount, ", slaveMessagesRecieved: ", slaveMessagesRecieved)
			slaveMessagesRecieved++
			// if slaveCount < 4 && slaveCount > -1 {
			// 	fmt.Println("Going to floor from slave: ", slaveCount)
			// 	// go goToFloor(slaveCount, &elevatorState)
			// }
			// sendImAliveMessage()
			time.Sleep(delay100ms) // wait 50 ms
			break
		case <-time.After(30 * delay100ms): // Wait 10 cycles (1 second). Master assumed dead
			// print("\033[H\033[2J")
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			// fmt.Println("Have not recieved slave message for the last 3 seconds! ", notRecievedCounter)
			notRecievedCounter++
			// newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			// err := newSlave.Run()
			// check(err)
			// return
		}
	}
}

func SendMasterIsAliveRegularly() {
	fmt.Println("Sending MasterIsAlive over network")

	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+masterIsAlivePort)
	check(err)
	msg := make([]byte, 8)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	defer func() {
		fmt.Println("Actually stopping sending MasterIsAliveMessages")
		udpBroadcast.Close()
	}()

	binary.BigEndian.PutUint64(msg, uint64(66))
	for {
		// fmt.Println("Sending I'm Alive from Master, msg:", msg)
		udpBroadcast.Write(msg)
		time.Sleep(delay100ms * 10)
	}
}

func ListenAfterAliveMasterRegularly(masterIsAliveChan chan int, stopListening chan bool) {
	fmt.Println("Listening to check if master is alive")
	udpAddr, err := net.ResolveUDPAddr("udp", masterIsAlivePort)
	check(err)
	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)
	defer udpListen.Close()

	// masterMessagesRecieved := 0

	buf := make([]byte, 8)

	go func() {
		for {
			// fmt.Println("Reading fromUDPbuf")
			udpListen.ReadFromUDP(buf)
			masterIsAliveChan <- int(binary.BigEndian.Uint64(buf))
			time.Sleep(100 * time.Millisecond) // Wait 1 cycle (100 ms)
		}

	}()

	for {
		// Convert byte from buf to int and send over channel.
		select {
		case <-stopListening:
			return
		}
		// masterIsAliveChan <- int(binary.BigEndian.Uint64(buf))
	}

}

func CheckIfMasterAlreadyExist() bool {
	fmt.Print("Are there any Masters here? ")
	udpAddr, err := net.ResolveUDPAddr("udp", masterIsAlivePort)
	// fmt.Println("updAddr: ", udpAddr)
	check(err)
	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)
	defer udpListen.Close()

	listenChan := make(chan int)

	// masterMessagesRecieved := 0

	go func() {
		buf := make([]byte, 8)
		for {
			udpListen.ReadFromUDP(buf)
			// Convert byte from buf to int and send over channel.
			listenChan <- int(binary.BigEndian.Uint64(buf))
			fmt.Println("Got message that master already exist, buf:", buf)
			time.Sleep(delay100ms) // Wait 1 cycle (100 ms)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	for {
		select {
		case <-listenChan:
			// udpListen.Close() /*Close instead*/
			time.Sleep(delay100ms) // wait 50 ms
			fmt.Println("YEEEES!")
			return true
		case <-time.After(4 * time.Second): // Wait 10 cycles (1 second). Master assumed dead
			// udpListen.Close() /*Close instead*/
			time.Sleep(time.Second)
			fmt.Println("NOOOO!")
			return false
		}
	}

}

func setOrderOverNetwork(destinationFloor int) {
	fmt.Println("Sending order over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+port)
	check(err)
	msg := make([]byte, 8)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	binary.BigEndian.PutUint64(msg, uint64(destinationFloor))
	udpBroadcast.Write(msg)
}

func SendJSON(m interface{}) {
	defer fmt.Println("Finished sending JSON")

	fmt.Println("Sending JSON over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+jsonSendPort)
	check(err)
	// msg := make([]byte, 128)
	// m := definitions.TestMessage{"Alice", "Hello", 1294706395881547000}
	b, _ := json.Marshal(m)

	// b :=

	fmt.Println("JSON in ByteArray:", b)
	jsonByteLength := len(b)
	firstByte := jsonByteLength / 255
	// fmt.Println("firstByte", firstByte)
	secondByte := jsonByteLength - firstByte*255
	// fmt.Println("secondByte:",secondByte)
	// fmt.Println("JSONByteArrayLength:",jsonByteLength)

	fmt.Println(byte(len(b)))

	b = append([]byte{byte(secondByte)}, b...)
	b = append([]byte{byte(firstByte)}, b...)
	// fmt.Println("WITH Length as first byte", b)

	// Create bcast Conn
	udpBroadcast, _ := net.DialUDP("udp", nil, udpAddr)
	// check(err)

	udpBroadcast.Write(b)
}

func RecieveJSON(updatedOrderList chan definitions.Orders) {
	fmt.Println("Listening after JSON Objectes")
	udpAddr, err := net.ResolveUDPAddr("udp", jsonSendPort)
	check(err)

	JSONobjectsRecieved := 0

	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)
	m := definitions.Orders{}

	listenChan := make(chan definitions.Orders, 1)
	// slaveCount := 0

	go func() {
		buf := make([]byte, 65536) /*2^16 = max recovery size*/
		for {
			udpListen.ReadFromUDP(buf)

			// fmt.Println("buffer after read from UDP: ", buf)

			// Two first bytes contains the size of the JSON byte array
			jsonByteLength := int(buf[0])*255 + int(buf[1])
			// fmt.Println("jsonByteLength:",jsonByteLength)

			// Convert byte from buf to int and send over channel.
			err := json.Unmarshal(buf[2:jsonByteLength+2], &m)
			// fmt.Println("Her kommer m som du skal se på: ", m)
			// fmt.Println("Ferdig med å vise m")
			check(err)
			listenChan <- m
			time.Sleep(delay100ms)
		}
	}()

	for {
		select {
		case JSONByteArray := <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("got JSON object: ", JSONByteArray, ", json objects recieved: ", JSONobjectsRecieved)
			storage.SaveJSONtoFile(JSONByteArray) //This actually works
			JSONobjectsRecieved++
			updatedOrderList <- JSONByteArray
			// if slaveCount < 4 && slaveCount > -1 {
			// 	fmt.Println("Going to floor from slave: ", slaveCount)
			// 	// go goToFloor(slaveCount, &elevatorState)
			// }
			// sendImAliveMessage()
			time.Sleep(delay100ms / 2) // wait 50 ms
			break
		case <-time.After(30 * delay100ms): // Wait 10 cycles (1 second). Master assumed dead
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			// fmt.Println("Have not recieved any JSON message for the last 3 seconds!")
			// newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			// err := newSlave.Run()
			check(err)
			// return
		}
	}

}

// func CheckForOrderListUpdatesFromMaster() {

// }

// var localIP string

// func getLocalIP() (string, error) {
// 	if localIP == "" {
// 		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
// 		if err != nil {
// 			return "", err
// 		}
// 		defer conn.Close()
// 		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
// 	}
// 	return localIP, nil
// }

// func CheckIfMasterAlreadyExist() {

// // }
// func ListenAfterAliveSlavesRegularly(aliveSlavesList *[]int) {
//   fmt.Println("Listening after ImAlive messages, list: ", aliveSlavesList)
//   udpAddr, err := net.ResolveUDPAddr("udp", slaveIsAlivePort)
//   check(err)

//   slaveMessagesRecieved := 0

//   // Create listen Conn
//   udpListen, err := net.ListenUDP("udp", udpAddr)
//   check(err)

//   listenChan := make(chan int, 1)
//   slaveCount := 0

//   go func() {
//     buf := make([]byte, 8)
//     for {
//       udpListen.ReadFromUDP(buf)

//       // Convert byte from buf to int and send over channel.
//       listenChan <- int(binary.BigEndian.Uint64(buf))
//       time.Sleep(delay100ms) // Wait 1 cycle (100 ms)
//     }
//   }()
//   notRecievedCounter := 0
//   for {
//     select {
//     case slaveCount = <-listenChan:
//       // fmt.Println("slaveCount: ", slaveCount)
//       fmt.Println("ListenAfterSlaves: ", slaveCount, ", slaveMessagesRecieved: ", slaveMessagesRecieved)
//       slaveMessagesRecieved++
//       // if slaveCount < 4 && slaveCount > -1 {
//       //  fmt.Println("Going to floor from slave: ", slaveCount)
//       //  // go goToFloor(slaveCount, &elevatorState)
//       // }
//       // sendImAliveMessage()
//       time.Sleep(delay100ms) // wait 50 ms
//       break
//     case <-time.After(30 * delay100ms): // Wait 10 cycles (1 second). Master assumed dead
//       // print("\033[H\033[2J")
//       // When master dies, slavecount is returned so that a new process of master -> slave
//       // can continue from the last value sent over the network.
//       fmt.Println("Have not recieved slave message for the last 3 seconds! ", notRecievedCounter)
//       notRecievedCounter++
//       // newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
//       // err := newSlave.Run()
//       // check(err)
//       // return
//     }
//   }

// }

// func decodeJSONRecievedAndStore(){

// }

// func SendFromSlaveToMaster(externalButtonPressesChan chan, ordersChan chan, elevatorState definitions.ElevatorState, id string) {
//   MSG := definitions.MSG_to_master {
//     Orders: Orders,
//     ElevatorState: state,
//     ExternalButtonPresses: <-externalButtonPressesChan,
//     Id: id,//Id string
//   }

//   sendJSON(MSG)
// }

// func SendFromMasterToSlave(elevators){
//   MSG := definitions.MSG_to_slave {
//     Elevators: elevators,
//   }
//   sendJSON(MSG)
// }

// func ListenToMaster() {

// }

// func ListenToSlave() {

// }
