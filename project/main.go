package main

import (
	//"controller"
	//"master"
	"./src/driver"
	//"./src/network"
	"fmt"
	//"os"
	//"time"
)

func main() {
	fmt.Println("Test")
	//network.Run()

	driver.Elev_init();
	driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)

	//driver.Elev_set_motor_direction(driver.DIRECTION_UP);

	for true {
		fmt.Println(driver.Elev_get_floor_sensor_signal())
		fmt.Println(driver.DIRECTION_DOWN)

		//driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)

		if driver.Elev_get_floor_sensor_signal() == driver.N_FLOORS - 1 {
			fmt.Println("Bobby Brown")
			driver.Elev_set_motor_direction(driver.DIRECTION_DOWN)
		} else if driver.Elev_get_floor_sensor_signal() == 0{
			fmt.Println("Bobby Brown inverse")
			driver.Elev_set_motor_direction(driver.DIRECTION_UP)
		}

		if driver.Elev_get_stop_signal() != 0 {
			driver.Elev_set_motor_direction(driver.DIRECTION_STOP)
			return
		}
	}
}
