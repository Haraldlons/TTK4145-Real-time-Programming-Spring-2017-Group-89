package driver

import "fmt"

const (
	MOTOR_SPEED      = 2800
	N_FLOORS         = 4
	N_BUTTONS        = 3
	DIRECTION_DOWN   = -1
	DIRECTION_STOP   = 0
	DIRECTION_UP     = 1
	BUTTON_CALL_UP   = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND   = 2
)

var (
	lamp_channel_matrix [N_FLOORS][N_BUTTONS]int = [N_FLOORS][N_BUTTONS]int{
		{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
		{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
		{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
		{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	}
)

var (
	button_channel_matrix [N_FLOORS][N_BUTTONS]int = [N_FLOORS][N_BUTTONS]int{
		{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
		{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}
)

func Elev_init() int {
	//Initialize HW
	if Io_init() == 0 {
		return 0
	}

	//Turn off all lights
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			Elev_set_button_lamp(b, f, 0)
		}
	}
	Elev_set_stop_lamp(0)
	Elev_set_floor_indicator(0)
	Elev_set_floor_indicator(0)

	return 1
}

// Direction can be -1, 0, 1
func Elev_set_motor_direction(direction int) {
	if direction == 0 {
		Io_write_analog(MOTOR, 0)
	} else if direction > 0 {
		fmt.Println("Direction Up")
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, MOTOR_SPEED)
	} else if direction < 0 {
		fmt.Println("Direction Down")
		fmt.Printf("Motor direction: %d\n", MOTORDIR)
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func Elev_set_button_lamp(button int, floor int, value int) int {
	if floor < 0 || floor >= N_FLOORS {
		return -1
	}
	//More fault handling here probs

	if value != 0 {
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}

	return 1
}

func Elev_set_floor_indicator(floor int) int {
	if floor < 0 || floor >= N_FLOORS {
		return -1
	}

	//Binary encoding. One light must always be on (??????Petter notes)
	if (floor & 0x02) != 0 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if (floor & 0x01) != 0 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
	return 1
}

//Turn on stop-light if value != 0
func Elev_set_stop_lamp(value int) {
	if value != 0 {
		Io_set_bit(LIGHT_STOP)
	} else {
		Io_clear_bit(LIGHT_STOP)
	}
}

//Turn on "door open"-light if value != 0
func Elev_set_door_open_lamp(value int) {
	if value != 0 {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

//-----------------------------------------------------------------\\

func Elev_get_button_signal(button int, floor int) int {
	if floor < 0 || floor >= N_FLOORS {
		return -1
	}
	return Io_read_bit(button_channel_matrix[floor][button])
}

func Elev_get_floor_sensor_signal() int {
	if Io_read_bit(SENSOR_FLOOR1) != 0 {
		return 0
	} else if Io_read_bit(SENSOR_FLOOR2) != 0 {
		return 1
	} else if Io_read_bit(SENSOR_FLOOR3) != 0 {
		return 2
	} else if Io_read_bit(SENSOR_FLOOR4) != 0 {
		return 3
	} else {
		return -1
	}
}

func Elev_get_stop_signal() int {
	return Io_read_bit(STOP)
}

func Elev_get_obstruction_signal() int {
	return Io_read_bit(OBSTRUCTION)
}
