package network

import (
	//"bufio"
	"fmt"
	"net"
	"os"
	//"strconv"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

//var baddr *net.UDPAddr

func Test() {
	//Prepare address at any address at port 20020
	serverAddr, err := net.ResolveUDPAddr("udp4", ":20020")
	checkError(err)

	//Listen at selected port
	serverConn, err := net.ListenUDP("udp4", serverAddr)
	checkError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
