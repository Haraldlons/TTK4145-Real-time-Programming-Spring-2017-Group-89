package buttons

import (
	"../definitions"
	"../driver"
	// "fmt"
	"time"
)

func Check_button_internal(buttonPressesChan chan [definitions.N_FLOORS]int) {
	pressed := driver.Elev_get_button_signal(2, 0)
	var buttonArray [definitions.N_FLOORS]int
	// var lastButtonArray [definitions.N_FLOORS]int

	// fmt.Println("buttonArray: ", buttonArray)
	// fmt.Println("Check_button_internal started. Button pressed: ", pressed)
	// last_i := 99 /*Random number above N_FLOORS*/
	for {
		for i := 0; i < definitions.N_FLOORS; i++ {
			pressed = driver.Elev_get_button_signal(2, i)

			// fmt.Println("What is pressed: ", i , ": ", pressed)
			if pressed != 0 {
				// lastButtonArray[last_i] = 1
				buttonArray[i] = 1
				// if buttonArray != lastButtonArray {
				// lastButtonArray[last_i] = 0
				buttonPressesChan <- buttonArray
				time.Sleep(300 * time.Millisecond)
				// }
				// last_i = i
				// fmt.Println("Button pressed: " ,i, ": ", pressed)
				buttonArray[i] = 0
			}
			time.Sleep(time.Millisecond * 5)
		}
		time.Sleep(time.Millisecond * 10)
	}
	return
}

// func Check_button_external(button int) bool {
// 	press := driver.Elev_get_button_signal(int, int)
// }
// PLEASE; FOR GODS SAKE: DONT LOOK AT THE CODE!!
func Check_button_external(buttonPressesChan chan [definitions.N_FLOORS][2]int) {
	pressed := driver.Elev_get_button_signal(1, 0)
	var buttonArray [definitions.N_FLOORS][2]int
	// fmt.Println("buttonArray: ", buttonArray)
	// fmt.Println("Check_button_internal started. Button pressed: ", pressed)
	var lastButtonArray [definitions.N_FLOORS][2]int
	last_i := 0 /*Random number above N_FLOORS*/
	last_j := 0
	isFirstTime := true
	for {
		for i := 0; i < definitions.N_FLOORS; i++ {
			for j := 0; j < 2; j++ {
				pressed = driver.Elev_get_button_signal(j, i)
				// fmt.Println("What is pressed: ", i , ": ", pressed)
				if pressed != 0 {
					if !isFirstTime {
						lastButtonArray[last_i][last_j] = 1
					}
					buttonArray[i][j] = 1
					if buttonArray != lastButtonArray {
						// fmt.Println("lastButtonArray:", lastButtonArray)
						lastButtonArray[last_i][last_j] = 0
						buttonPressesChan <- buttonArray
						time.Sleep(200 * time.Millisecond)
					}
					last_i = i
					last_j = j
					// fmt.Println("Button pressed: " ,i, ": ", pressed)
					buttonArray[i][j] = 0
					isFirstTime = false
				}
				time.Sleep(time.Millisecond * 5)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
	return
}
