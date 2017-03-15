package def

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

	//Extra states
	IDLE      int = -1 // Used in the destination field of "Order"
	// MOVING    int = 1
	// DOOR_OPEN int = 

	BcAddress string = "129.241.187.255"
	MasterIsAlivePort string = ":46721"
	SlaveIsAlivePort string = ":46720"
	MasterToSlavePort string = ":18900"
	SlaveToMasterPort string = ":18901"	
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

type Elevators struct {
	OrderMap         map[string]Orders
	ElevatorStateMap map[string]ElevatorState
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
