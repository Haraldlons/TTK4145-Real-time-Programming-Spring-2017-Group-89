package buttons

import (
	"../def"
	"../driver"
	// "fmt"
	"time"
)

func Check_button_internal(buttonPressesChan chan [def.N_FLOORS]int) {
	pressed := driver.Elev_get_button_signal(2, 0)
	var buttonArray [def.N_FLOORS]int
	for {

		for i := 0; i < def.N_FLOORS; i++ {
			pressed = driver.Elev_get_button_signal(2, i)
			if pressed != 0 {
				buttonArray[i] = 1
					buttonPressesChan <- buttonArray
					time.Sleep(200 * time.Millisecond)
				buttonArray[i] = 0
			}
			time.Sleep(time.Millisecond * 5)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func Check_button_external(buttonPressesChan chan [def.N_FLOORS][2]int) {
	pressed := driver.Elev_get_button_signal(1, 0)
	var buttonArray [def.N_FLOORS][2]int
	for {
		for i := 0; i < def.N_FLOORS; i++ {
			for j := 0; j < 2; j++ {
				pressed = driver.Elev_get_button_signal(j, i)
				if pressed != 0 {
					buttonArray[i][j] = 1
						buttonPressesChan <- buttonArray
						time.Sleep(200 * time.Millisecond)
					buttonArray[i][j] = 0
				}
				time.Sleep(time.Millisecond * 5)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}
