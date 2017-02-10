package main


import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
)

const(
	hostaddr = "129.241.187.255:20020"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func udpListen(port string) {
	buf := make([]byte, 1024)

	//Prepare address
	serverAddr, err := net.ResolveUDPAddr("udp4", ":" + port)
	CheckError(err)

	//Listen at selected port
	serverConn, err := net.ListenUDP("udp4", serverAddr)
	CheckError(err)
	defer serverConn.Close()

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
            fmt.Println("Error: ",err)
        } 
	}
}

func udpSend() {
	ServerAddr, err := net.ResolveUDPAddr("udp4", hostaddr)
    CheckError(err)
 
    LocalAddr, err := net.ResolveUDPAddr("udp4", hostaddr)
    CheckError(err)
 
    Conn, err := net.DialUDP("udp4", LocalAddr, ServerAddr)
    CheckError(err)
 
    defer Conn.Close()
    i := 0
    for {
        msg := "Test" + strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _,err := Conn.Write(buf)

        if err != nil {
            fmt.Println(msg, err)
        }
        
        time.Sleep(time.Second * 1)
    }
}


func main() {
	udpSend()
	udpListen("20020")
}



