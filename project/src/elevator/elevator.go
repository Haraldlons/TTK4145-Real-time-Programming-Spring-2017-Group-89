package elevator

import (
	"fmt"
	"driver"
)

func ExecuteOrders(orders struct, elevatorState *definitions.ElevatorState) bool {
	for _, elem := range orders {
		goToFloor(elem.destinationFloor, elevatorState)
	}
}

func goToFloor(destinationFloor int, elevatorState *definitions.ElevatorState) bool {
	fmt.Println("Going to floor: ", destinationFloor)
	lastFloor := elevatorState.LastFloor

	if(driver.Elev_get_floor_sensor_signal() == destinationFloor){
		fmt.Println("You are already on the desired floor")
		driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
		return true
	} 
	else {  /*You are not on the desired floor*/
		if lastFloor < destinationFloor {
			driver.Elev_set_motor_direction(driver.DIRECTION_UP)
		} else {
			driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
		}
		for {
			if driver.Elev_get_floor_sensor_signal() == destinationFloor {
				fmt.Println("You reached your desired floor. Walk out")
				driver.Elev_set_floor_indicator(destinationFloor) // Turn on corresponding floor indicator
				driver.Elev_set_motor_direction(driver.DIRECTION_STOP) // Stop elevator
				return true
			}else {
				time.Sleep(delay)
			}
		}
	}


}

func GetState() struct {
	return state
}
