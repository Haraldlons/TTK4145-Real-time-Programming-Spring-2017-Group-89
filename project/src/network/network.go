package network

import (
	// "../definitions"
	// "../driver"
	// "./../controller"
	// "./src/network"
	// "../buttons"
	//"./src/driver"
	// "../storage"
	//"./src/master"
	// "../elevator"

	//"./src/watchdog"
	// "elevator"
	// "network"
	// "storage"
	"fmt"
	"time"
	 "encoding/binary"
// 	// "fmt"
// 	"net"
	"os/exec"

	//"bufio"
	// "log"
	"net"
	//"os"
	// "time"
	//"strconv"
)

// var listenTimer = 100 * time.Millisecond
// var broadcastTimer = 100 * time.Millisecond

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// // Function for checking arguments passed to network functions
// func checkArgs() {

// }

// func Broadcast(port string, bcAddr string, sendChan chan string) {
// 	addr, err := net.ResolveUDPAddr("udp", bcAddr+port)
// 	check(err)

// 	conn, err := net.DialUDP("udp", nil, addr)
// 	check(err)

// 	log.Println("Broadcasting connection established on: ", bcAddr+port)

// 	msg := make([]byte, 1024)
// 	for {
// 		conn.Write(msg)
// 		time.Sleep(broadcastTimer)
// 	}
// }

// func Listen(port string, listenChan chan string) {
// 	addr, err := net.ResolveUDPAddr("udp", port)
// 	check(err)

// 	conn, err := net.ListenUDP("udp", addr)
// 	check(err)

// 	log.Println("Listening connection established on port ", port)

// 	buf := make([]byte, 1024)
// 	for {
// 		conn.ReadFromUDP(buf)
// 		listenChan <- string(buf) // Dunno
// 		time.Sleep(listenTimer)
// 	}
// }
/*Above is the original code before my expermintal coding-session*/

// Experimental work done
// import (
// 	"encoding/binary"
// 	// "fmt"
// 	"net"
// 	"os/exec"
// 	// "time"
// )

var bcAddress string = "129.241.187.255"
var port string = ":55555"
var slaveSendPort string = ":55758"

var delay100ms = 100 * time.Millisecond

func check(err error) {
	if err != nil {
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
	go sendImAliveMessage(stopSlaveBroadcasting)

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

	go listenAfterImAliveMessage()

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
		time.Sleep(10*delay100ms) // Wait 1 cycle (100 ms)
	}
}

func SetupNetwork() {

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

func sendImAliveMessage(stopSlaveBroadcasting chan int){
	defer fmt.Println("Actually stopping sending sendImAliveMessage")

	fmt.Println("Sending order over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+slaveSendPort)
	check(err)
	msg := make([]byte, 8)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	binary.BigEndian.PutUint64(msg, uint64(66))
	// go func() {
	// 	select {
	// 	case <- stopSlaveBroadcasting:
	// 		fmt.Println("Ending sendImAliveMessage")
	// 		return
	// 	}
	// }()
	for {
		select {
			case <- stopSlaveBroadcasting:
				return
			default:
				udpBroadcast.Write(msg)
				// fmt.Println("Sending I'm Alive")
				time.Sleep(delay100ms*2)
		}
	}


}

func listenAfterImAliveMessage(){
		fmt.Println("Listening after ImAlive messages")
		udpAddr, err := net.ResolveUDPAddr("udp", slaveSendPort)
		check(err)

		slaveMessagesRecieved := 0

		// Create listen Conn
		udpListen, err := net.ListenUDP("udp", udpAddr)
		check(err)

		listenChan := make(chan int, 1)
		slaveCount := 0


		go func(){
			buf := make([]byte, 8)
			for {
				udpListen.ReadFromUDP(buf)

				// Convert byte from buf to int and send over channel.
				listenChan <- int(binary.BigEndian.Uint64(buf))
				time.Sleep(delay100ms) // Wait 1 cycle (100 ms)
			}
		}()



	for {
		select {
		case slaveCount = <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("ListenAfterSlaves: ", slaveCount, ", slaveMessagesRecieved: ", slaveMessagesRecieved)
			slaveMessagesRecieved++
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
			fmt.Println("Have not recieved slave message for the last 3 seconds!")
			newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			err := newSlave.Run()
			check(err)
			// return
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
