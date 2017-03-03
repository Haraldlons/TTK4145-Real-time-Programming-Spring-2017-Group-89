package watchdog

import (
	"net"
	"network"
	"time"
)

var timeLimit = 10 * time.Second
var listenTimer = 100 * time.Millisecond

func SendNetworkAlive(udpBroadcast *net.UDPConn) bool { // Possibly unnesces
	msg := make([]byte, 1024)
	numLines, err := udpBroadcast.Write(msg)
	return err != nil
}

func CheckNetworkAlive(udpListen *net.UDPConn) int {
	listenChan := make(chan int, 1)
	lifeCheck := 1

	// Run listening goroutine
	go listen(listenChan, udpListen)

	for {
		select {
		case lifeCheck = <-listenChan:
			if lifeCheck == 1 {
				time.Sleep(listenTimer)
			} else { // Possibly dangerous. RETHINK!
				return -1
			}
		case <-time.After(timeLimit): // Node assumed dead
			return -1
		}
	}
}

func listen(listenChan chan int, udpListen *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		udpListen.ReadFromUDP(buf)
		listenChan <- int(buf)
		time.Sleep(listenTimer)
	}
}

/*
func CheckElevatorState(state var) var {
	state := true
	return state
}

func reset_Master() {

}

func reset_Elevator() {

}
*/
