package definitions

// import "fmt"

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
	ON               bool = true
	OFF              bool = false

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
	Direction int
}

type Orders struct {
	Orders []Order
}

type ElevatorState struct {
	LastFloor   int
	Direction   int
	Destination int
}

type ElevatorStateMap struct {
	map [string] ElevatorState
}

type Elevators struct {
	OrderMap       map[string]Orders
	ElevatorStateMap ElevatorStateMap
}

type MSG_to_master struct {
	Orders                Orders
	ElevatorState         ElevatorState
	ExternalButtonPresses []Order
	Id                    string
}

type MSG_to_slave struct {
	Elevators Elevators
}

type TestMessage struct {
	Name string
	Body string
	Time int64
}
