package main

import (
	"fmt"
	"net"
	"time"
	"encoding/binary"
)

var address string = "localhost"
var bcAddress string = "129.241.187.255"
var port string = "20020"
var delay = 100*time.Millisecond

func check(err error){
	if err != nil {
		panic(err)
	}
}

func slave(udpListen *net.UDPConn) int {
	listenChan := make(chan int, 1)
	slaveCount := 0
	go listen(listenChan, udpListen)
	for {	
		select {
		case slaveCount <- listenChan:
			time.Sleep(delay/2) // wait 50 ms 
			break
		case <- time.After(10*delay) // Wait 10 cycles (1 second)
			fmt.Println("Mufasa is dead. Long live the king.")
			return slaveCount
 		}
	}
}

func master(startCount int, udpBroadcast *net.UDPconn) {
	
	msg := make([]byte, 1)

	for {
		count++
		fmt.Println(count)
		msg[0] = byte(count) 
		udpBroadcast.Write(msg)
		time.Sleep(delay)
	}
}

func listen(listenChan chan int ,udpListen *net.UDPConn){
	buf := make(byte[], 1024)

	for {
		udpListen.ReadFromUDP(buf)

		// Convert byte from buf to int and send over channel.
		listenChan <- binary.BigEndian.Uint64(buf) 
		time.Sleep(delay)
	}
}

func changeMaster(){
	fmt.Println("Changing master")
}

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", port)
	check(err)

	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)

	// Initialize slave
	count := slave(udpListen, )

	udpListen.Close()

	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress + ":" + port)
	check(err)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	fmt.Println("Run primary")
	master(count, udpBroadcast)

	fmt.Println("Close connections")
	udpBroadcast.Close()
}
