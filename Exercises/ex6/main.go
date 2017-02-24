package main

import (
	"fmt"
	"net"
)

var address string = "localhost"
var bcAdress string = "129.241.187.255"
var port string = "20020"

func check(err error){
	if err != nil {
		panic(err)
	}
}

func slave(conn net.UDPConn, isAlive bool) bool {

	fmt.Println("Mufasa is dead. Long live the king.")
}

func master(startCount int, udpBroadcast *net.UDPconn) {


	msg := make([]byte, 1)
	count++
	fmt.Println(*count)
}

func sendAlive(){

}

func listenAlive(){

}

func changeMaster(){

}

func main() {
	fmt.Println("test")
	isAlive bool := false
	count := 0

	udpAddr, err := net.ResolveUDPAddr("udp", port)
	check(err)

	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	fmt.Println("Run primary")
	master(count, udpBroadcast)

	defer udpBroadcast.close
}
