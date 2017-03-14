package network

import (
	"../definitions"
	// "../driver"
	// "../buttons"
	"../storage"
	// "../elevator"
	// "encoding/binary"
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
	"bytes"
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

func SendSlaveIsAliveRegularly(elevator_id string, stopSending chan bool) {
	// fmt.Println("Sending ImAliveMessage over network")

	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+slaveIsAlivePort)
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)

	defer func() {
		fmt.Println("Actually stopping sending sendImAliveMessage")
		udpBroadcast.Close()
	}()

	if err != nil { //Can't connect to the interwebs
		udpAddr, _ = net.ResolveUDPAddr("udp", "localhost"+masterIsAlivePort)
		udpBroadcast, _ = net.DialUDP("udp", nil, udpAddr)
	}

	msg := []byte(elevator_id)
	for {
		select {
		case <-stopSending:
			// time.Sleep(10 * time.Millisecond)
			return
		default:
			udpBroadcast.Write(msg)
			// fmt.Println("Sending I'm Alive")
			time.Sleep(time.Millisecond * 200)
		}
	}
}

func ListenAfterAliveSlavesRegularly(updatedSlaveIdChanMap map[string]chan string, stopListening bool) {
	// allSlavesMap := make(map[string]bool)

	// Create listen Conn
	udpAddr, _ := net.ResolveUDPAddr("udp", slaveIsAlivePort)
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	defer udpListen.Close()

	go func() {
		buf := make([]byte, 16)
		for {
			udpListen.ReadFromUDP(buf)

			// Convert buf from byte to string (IP-address)
			n := bytes.IndexByte(buf, 0)
			slave_id := string(buf[:n])

			// Send update to "keepTrackOfAllAliveSlaves()" and watchDog
			updatedSlaveIdChanMap["toKeepTrackOfAllAliveSlaves"] <- slave_id
			updatedSlaveIdChanMap["toWatchdog"] <- slave_id

			time.Sleep(time.Millisecond * 100) // Wait 1 cycle (100 ms)
		}
	}()

	for { //Infinite loop
		time.Sleep(time.Second)
	}
}

func SendMasterIsAliveRegularly(master_id string, stopSending chan bool) {
	udpAddr, _ := net.ResolveUDPAddr("udp", bcAddress+masterIsAlivePort)
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil { //Can't connect to the interwebs
		udpAddr, _ = net.ResolveUDPAddr("udp", "localhost"+masterIsAlivePort)
		udpBroadcast, err = net.DialUDP("udp", nil, udpAddr)
	}

	defer func() {
		fmt.Println("Actually stopping sending MasterIsAliveMessages")
		udpBroadcast.Close()
	}()

	msg := []byte(master_id)
	for {
		select {
		case <-stopSending:
			return
		default:
			fmt.Println("Sending I'm Alive from Master, msg:", msg)
			udpBroadcast.Write(msg)
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func ListenAfterAliveMasterRegularly(masterIsAliveChan chan string, stopListening chan bool) {
	// fmt.Println("Listening to check if master is alive")
	udpAddr, _ := net.ResolveUDPAddr("udp", masterIsAlivePort)
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	defer udpListen.Close()

	buf := make([]byte, 16)
	go func() {
		for {
			udpListen.ReadFromUDP(buf)
			n := bytes.IndexByte(buf, 0)
			master_id := string(buf[:n])

			masterIsAliveChan <- master_id
			time.Sleep(100 * time.Millisecond) // Wait 1 cycle (100 ms)
		}

	}()

	for {
		select {
		case <-stopListening:
			return
		}
	}
}

func CheckIfMasterAlreadyExist() bool {
	fmt.Print("Are there any Masters here? ")

	udpAddr, _ := net.ResolveUDPAddr("udp", masterIsAlivePort)
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	defer udpListen.Close()

	listenChan := make(chan string)

	go func() {
		buf := make([]byte, 16)
		for {
			udpListen.ReadFromUDP(buf)
			n := bytes.IndexByte(buf, 0)
			master_id := string(buf[:n])

			listenChan <- master_id
			time.Sleep(delay100ms) // Wait 1 cycle (100 ms)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	for {
		select {
		case <-listenChan:
			fmt.Println("YEEEES!")
			return true
		case <-time.After(2 * time.Second): // Master assumed dead
			time.Sleep(time.Second)
			fmt.Println("NOOOO!")
			return false
		}
	}
}

func SendUpdatesToMaster(msg definitions.MSG_to_master, elevatorState definitions.ElevatorState, elevator_id string, lastSentMsgToMasterChanForPrinting chan<- definitions.MSG_to_master) {
	msg.Id = elevator_id
	msg.ElevatorState = elevatorState
	// defer fmt.Println("Finished sending JSON")
	fmt.Println("Sending UPDATED OrderList To Master. MSG_to_master: ", msg)
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+slaveToMasterPort)
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil { //Can't connect to the interwebs
		udpAddr, _ = net.ResolveUDPAddr("udp", "localhost"+slaveToMasterPort)
		udpBroadcast, _ = net.DialUDP("udp", nil, udpAddr)
	}

	defer udpBroadcast.Close()
	lastSentMsgToMasterChanForPrinting<- msg
	//check(_)
	// msg := make([]byte, 128)
	// msg := definitions.TestMessage{"Alice", "Hello", 1294706395881547000}
	b, _ := json.Marshal(msg)

	// b :=

	// fmt.Println("JSON in ByteArray:", b)
	jsonByteLength := len(b)
	firstByte := jsonByteLength / 255
	// fmt.Println("firstByte", firstByte)
	secondByte := jsonByteLength - firstByte*255
	// fmt.Println("secondByte:",secondByte)
	// fmt.Println("JSONByteArrayLength:",jsonByteLength)

	// fmt.Println(byte(len(b)))

	b = append([]byte{byte(secondByte)}, b...)
	b = append([]byte{byte(firstByte)}, b...)
	// fmt.Println("WITH Length as first byte", b)

	// Create bcast Conn
	// //check(_)

	udpBroadcast.Write(b)
	return
}

func ListenToMasterUpdates(updatedOrderList chan definitions.Orders, elevator_id string, lastRecievedMSGFromMasterChanForPrinting chan<- definitions.MSG_to_slave) {
	fmt.Println("Listening after Updates from Master")

	udpAddr, _ := net.ResolveUDPAddr("udp", masterToSlavePort)
	udpListen, _ := net.ListenUDP("udp", udpAddr)
	defer udpListen.Close()

	msg := definitions.MSG_to_slave{}

	listenChan := make(chan definitions.MSG_to_slave)

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
			// fmt.Println("Her kommer msg som du skal se på: ", msg)
			// fmt.Println("Ferdig med å vise msg")
			//check(_)
			listenChan <- msg
			lastRecievedMSGFromMasterChanForPrinting <- msg
			time.Sleep(delay100ms)
		}
	}()

	for {
		select {
		case MSG_to_slave := <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("Message received from master:", MSG_to_slave)
			storage.SaveJSONtoFile(MSG_to_slave.Elevators) //This actually works
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

func GetLocalIP() (string, error) {
	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
	if err != nil {
		return "localhost", err
	}
	defer conn.Close()
	localIP := strings.Split(conn.LocalAddr().String(), ":")[0]

	return localIP, nil
}

func ListenToSlave(msgChan chan definitions.MSG_to_master) {
	fmt.Println("Listening after messages from slave")
	udpAddr, _ := net.ResolveUDPAddr("udp", slaveToMasterPort)
	udpListen, _ := net.ListenUDP("udp", udpAddr)

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
				// Decode message
				json.Unmarshal(buf[2:jsonByteLength+2], &msg)

				msgChan <- msg
			}
			time.Sleep(10 * time.Millisecond)
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
