package network

// import (
// 	//"bufio"
// 	"log"
// 	"net"
// 	//"os"
// 	"time"
// 	//"strconv"
// )

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
