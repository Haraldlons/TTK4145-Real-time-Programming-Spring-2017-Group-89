package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"time"
)

var bcAddress string = "129.241.187.255"
var port string = ":55555"
var delay = 300 * time.Millisecond

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
			fmt.Println("slaveCount: ", slaveCount)
			time.Sleep(delay / 2) // wait 50 ms
			break
		case <-time.After(10 * delay): // Wait 10 cycles (1 second). Master assumed dead
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			fmt.Println("Mufasa is dead. Long live the king.")
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

	count := startCount
	msg := make([]byte, 8)

	for {
		// Convert count from int to binary/byte and place in msg
		binary.BigEndian.PutUint64(msg, uint64(count))
		udpBroadcast.Write(msg)

		fmt.Println(count)
		count++

		time.Sleep(delay) // Wait 1 cycle (100 ms)
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

func main() {

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

	fmt.Println("Run master")
	master(count, udpBroadcast)

	fmt.Println("Close broadcast")
	udpBroadcast.Close()
}
