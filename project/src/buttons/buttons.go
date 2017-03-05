package buttons

import (
	"../driver"
	// "../definitions"
	"time"
	"fmt"
)

func Check_button_internal(buttonPressesChan chan int) {
	pressed := driver.Elev_get_button_signal(2,0)
	fmt.Println("Check_button_internal started. Button pressed: ", pressed)
	for {
		for i := 0; i < 4; i++ {
			pressed = driver.Elev_get_button_signal(2,i)
			fmt.Println("What is pressed: ", i , ": ", pressed)
			// if pressed != 1{
			// 	fmt.Println("Button pressed: " ,i, ": ", pressed)
			// 	buttonPressesChan <- i
			// }
			time.Sleep(time.Millisecond*50)
		}
		time.Sleep(time.Millisecond * 10)
	}
	return
}

// func Check_button_external(button int) bool {
// 	press := driver.Elev_get_button_signal(int, int)
// }

func Check_button_external(buttonPressesChan chan int) {
	pressed := driver.Elev_get_button_signal(1,0)
	fmt.Println("Check_button_internal started. Button pressed: ", pressed)
	for {
		for i := 2; i < 3; i++ {
			pressed = driver.Elev_get_button_signal(1,i)
			fmt.Println("What is pressed: ", i , ": ", pressed)
			// if pressed != 1{
			// 	fmt.Println("Button pressed: " ,i, ": ", pressed)
			// 	buttonPressesChan <- i
			// }
			time.Sleep(time.Millisecond*50)
		}
		time.Sleep(time.Millisecond * 10)
	}
	return
}