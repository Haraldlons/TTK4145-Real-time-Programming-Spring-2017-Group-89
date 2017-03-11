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
	// fmt.Println("buttonArray: ", buttonArray)
	// fmt.Println("Check_button_internal started. Button pressed: ", pressed)
	for {
		for i := 0; i < definitions.N_FLOORS; i++ {
			pressed = driver.Elev_get_button_signal(2, i)

			// fmt.Println("What is pressed: ", i , ": ", pressed)
			if pressed != 0 {
				buttonArray[i] = 1
				// fmt.Println("Button pressed: " ,i, ": ", pressed)
				buttonPressesChan <- buttonArray
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

func Check_button_external(buttonPressesChan chan [definitions.N_FLOORS][2]int) {
	pressed := driver.Elev_get_button_signal(1, 0)
	var buttonArray [definitions.N_FLOORS][2]int
	// fmt.Println("buttonArray: ", buttonArray)
	// fmt.Println("Check_button_internal started. Button pressed: ", pressed)
	for {
		for i := 0; i < definitions.N_FLOORS; i++ {

			for j := 0; j < 2; j++ {
				pressed = driver.Elev_get_button_signal(j, i)
				// fmt.Println("What is pressed: ", i , ": ", pressed)
				if pressed != 0 {
					buttonArray[i][j] = 1
					// fmt.Println("Button pressed: " ,i, ": ", pressed)
					buttonPressesChan <- buttonArray
					buttonArray[i][j] = 0
				}
				time.Sleep(time.Millisecond * 5)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
	return
}
