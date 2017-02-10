package main

import(
	"./driver"
	"fmt"
)

int main() {
    driver.Elev_init();

    driver.Elev_set_motor_direction(DIRN_UP);

    for true {
        // Change direction when we reach top/bottom floor
        if driver.Elev_get_floor_sensor_signal() == N_FLOORS - 1 {
            driver.Elev_set_motor_direction(DIRN_DOWN);
        } else if driver.Elev_get_floor_sensor_signal() == 0 {
            driver.Elev_set_motor_direction(DIRN_UP);
        }

        // Stop elevator and exit program if the stop button is pressed
        if driver.Elev_get_stop_signal() {
            driver.Elev_set_motor_direction(DIRN_STOP);
            return 0;
        }
    }
}
