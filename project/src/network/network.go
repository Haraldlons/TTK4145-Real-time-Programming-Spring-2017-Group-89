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
	// "os/exec"
	"time"
	// "math"
	//"bufio"
	// "log"
	//"os"
	//"strconv"
	"strings"
)

var bcAddress string = "129.241.187.255"

//var bcAddress string = "localhost"
var port string = ":46723"
var slaveIsAlivePort string = ":46720"
var masterIsAlivePort string = ":46721"
var jsonSendPort string = ":46724"
var masterToSlavePort string = ":18900"
var slaveToMasterPort string = ":18901"

var delay100ms = 100 * time.Millisecond

func SendSlaveIsAliveRegularly(slaveID int, stopSlaveBroadcasting chan bool) {
	fmt.Println("Sending ImAliveMessage over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+slaveIsAlivePort)
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil { //Can't connect to the interwebs
		// fmt.Println("err is not nil", err)
		udpAddr, _ = net.ResolveUDPAddr("udp", "localhost"+masterIsAlivePort)
		udpBroadcast, _ = net.DialUDP("udp", nil, udpAddr)
	}

	//check(_)
	msg := make([]byte, 8)

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
	udpAddr, _ := net.ResolveUDPAddr("udp", slaveIsAlivePort)

	//check(_)

	slaveMessagesRecieved := 0

	// Create listen Conn
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	//check(_)
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
			// _ := newSlave.Run()
			// //check(_)
			// return
		}
	}
}

func SendMasterIsAliveRegularly() {
	fmt.Println("Sending MasterIsAlive over network")

	udpAddr, _ := net.ResolveUDPAddr("udp", bcAddress+masterIsAlivePort)
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil { //Can't connect to the interwebs
		udpAddr, _ = net.ResolveUDPAddr("udp", "localhost"+masterIsAlivePort)
		udpBroadcast, err = net.DialUDP("udp", nil, udpAddr)
	}
	fmt.Println("udpBroadcast: ", udpBroadcast)

	msg := make([]byte, 8)
	defer func() {
		fmt.Println("Actually stopping sending MasterIsAliveMessages")
		udpBroadcast.Close()
	}()

	binary.BigEndian.PutUint64(msg, uint64(66))
	for {
		fmt.Println("Sending I'm Alive from Master, msg:", msg)
		udpBroadcast.Write(msg)
		time.Sleep(delay100ms * 10)
	}
}

func ListenAfterAliveMasterRegularly(masterIsAliveChan chan int, stopListening chan bool) {
	fmt.Println("Listening to check if master is alive")
	udpAddr, _ := net.ResolveUDPAddr("udp", masterIsAlivePort)
	//check(_)
	// Create listen Conn
	udpListen, _ := net.ListenUDP("udp", udpAddr)

	//check(_)
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
	udpAddr, _ := net.ResolveUDPAddr("udp", masterIsAlivePort)
	// fmt.Println("updAddr: ", udpAddr)
	//check(_)
	// Create listen Conn
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	//check(_)
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
		case <-time.After(2 * time.Second): // Wait 10 cycles (1 second). Master assumed dead
			// udpListen.Close() /*Close instead*/
			time.Sleep(time.Second)
			fmt.Println("NOOOO!")
			return false
		}
	}

}

func SendJSON(m interface{}) {
	defer fmt.Println("Finished sending JSON")
	fmt.Println("Sending JSON over network. Interface: ", m)
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+slaveToMasterPort)
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil { //Can't connect to the interwebs
		udpAddr, _ = net.ResolveUDPAddr("udp", "localhost"+slaveToMasterPort)
		udpBroadcast, _ = net.DialUDP("udp", nil, udpAddr)
	}

	defer udpBroadcast.Close()
	//check(_)
	// msg := make([]byte, 128)
	// m := definitions.TestMessage{"Alice", "Hello", 1294706395881547000}
	b, _ := json.Marshal(m)

	// b :=

	// fmt.Println("JSON in ByteArray:", b)
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
	// //check(_)

	udpBroadcast.Write(b)
}

func ListenToMasterUpdates(updatedOrderList chan definitions.Orders, elevator_id string) {
	fmt.Println("Listening after Updates from Master")
	udpAddr, _ := net.ResolveUDPAddr("udp", masterToSlavePort)
	//check(_)
	// Create listen Conn
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	//check(_)
	msg := definitions.MSG_to_slave{}

	messagesRecievedFromMaster := 0
	listenChan := make(chan definitions.MSG_to_slave, 1)

	go func() {
		buf := make([]byte, 65536) /*2^16 = max recovery size*/
		for {
			udpListen.ReadFromUDP(buf)

			// fmt.Println("buffer after read from UDP: ", buf)

			// Two first bytes contains the size of the JSON byte array
			jsonByteLength := int(buf[0])*255 + int(buf[1])
			// fmt.Println("jsonByteLength:",jsonByteLength)

			// Convert byte from buf to int and send over channel.
			json.Unmarshal(buf[2:jsonByteLength+2], &msg)
			// fmt.Println("Her kommer m som du skal se på: ", m)
			// fmt.Println("Ferdig med å vise m")
			//check(_)
			listenChan <- msg
			time.Sleep(delay100ms)
		}
	}()

	for {
		select {
		case MSG_to_slave := <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("got MSG_to_slave object: ", MSG_to_slave, ", json objects recieved: ", messagesRecievedFromMaster)
			storage.SaveJSONtoFile(MSG_to_slave.Elevators) //This actually works
			messagesRecievedFromMaster++
			updatedOrderList <- MSG_to_slave.Elevators.OrderMap[elevator_id]
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
			// _ := newSlave.Run()
			//check(_)
			// return
		}
	}

}

// func CheckForOrderListUpdatesFromMaster() {

// }

var localIP string

func GetLocalIP() (string, error) {
	fmt.Println("Checking IP")
	if localIP == "" {
		fmt.Println("Testing")
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	fmt.Println("Our local IP: ", localIP)
	return localIP, nil
}

// func CheckIfMasterAlreadyExist() {

// // }
// func ListenAfterAliveSlavesRegularly(aliveSlavesList *[]int) {
//   fmt.Println("Listening after ImAlive messages, list: ", aliveSlavesList)
//   udpAddr, _ := net.ResolveUDPAddr("udp", slaveIsAlivePort)
//   //check(_)

//   slaveMessagesRecieved := 0

//   // Create listen Conn
//   udpListen, _ := net.ListenUDP("udp", udpAddr)
//   //check(_)

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
//       // _ := newSlave.Run()
//       // //check(_)
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

func ListenToSlave(msgChan chan definitions.MSG_to_master) {
	fmt.Println("Listening after messages from slave")
	udpAddr, _ := net.ResolveUDPAddr("udp", slaveToMasterPort)
	//check(_)

	// Create listen Conn
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	fmt.Println("udpListen:", *udpListen)
	//check(_)
	defer udpListen.Close()

	go func() {
		// Buffer for received message
		buf := make([]byte, 65536) /*2^16 = max recovery size*/
		for {
			// Listen for messages
			udpListen.ReadFromUDP(buf)
			// Two first bytes contains the size of the JSON byte array
			jsonByteLength := int(buf[0])*255 + int(buf[1])
			msg := definitions.MSG_to_master{}
			// Convert back to struct
			if jsonByteLength > 0 {
				// fmt.Println("length of json:", jsonByteLength)
				json.Unmarshal(buf[2:jsonByteLength+2], &msg)
				// fmt.Println("Recieved json object from slave. ", msg)
				//check(_)
				// fmt.Println("after error")
				msgChan <- msg
			}
			time.Sleep(1000 * time.Millisecond)
			// Send message over channel
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}

func SendToSlave(msg definitions.MSG_to_slave /*udpBroadcast *net.UDPConn*/) {

	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+masterToSlavePort)
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil { //Can't connect to the interwebs
		fmt.Println("err is not nil", err)
		udpAddr, err = net.ResolveUDPAddr("udp", "localhost"+masterToSlavePort)
		udpBroadcast, err = net.DialUDP("udp", nil, udpAddr)
	}
	//check(_)

	defer udpBroadcast.Close()
	buf, _ := json.Marshal(msg)
	// fmt.Println("JSON in ByteArray:", buf)
	jsonByteLength := len(buf)
	firstByte := jsonByteLength / 255
	// fmt.Println("firstByte", firstByte)
	secondByte := jsonByteLength - firstByte*255
	// fmt.Println("secondByte:",secondByte)
	// fmt.Println("JSONByteArrayLength:",jsonByteLength)

	// fmt.Println(byte(len(buf)))

	buf = append([]byte{byte(secondByte)}, buf...)
	buf = append([]byte{byte(firstByte)}, buf...)

	udpBroadcast.Write(buf)
	// defer fmt.Println("Have sent message to Slave, buf: ", buf)
	defer fmt.Println("Actual message: ", msg)
	return
}
