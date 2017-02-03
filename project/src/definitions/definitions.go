package definitions

import "fmt"

const (
	//General constants
	N_ELEVS          int  = 3
	N_FLOORS         int  = 4
	BUTTON_CALL_UP   int  = 0
	BUTTON_CALL_DOWN int  = 1
	BUTTON_COMMAND   int  = 2 //USED FOR ????
	DIR_UP           int  = 1
	DIR_DOWN         int  = -1
	DIR_STOP         int  = 0
	ON               bool = 1
	OFF              bool = 0

	//States
	IDLE      int = 0
	MOVING    int = 1
	DOOR_OPEN int = 2

	LOCAL_LISTEN_PORT     int = 20020
	BROADCAST_LISTEN_PORT int = 30020
	MESSAGE_SIZE          int = 1024
)

type Order struct {
	Floor     int
	Direction bool
}

type MSG struct {
	State int
	//Add more when more flushed out
}
