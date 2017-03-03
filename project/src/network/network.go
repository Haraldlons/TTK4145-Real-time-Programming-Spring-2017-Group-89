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
