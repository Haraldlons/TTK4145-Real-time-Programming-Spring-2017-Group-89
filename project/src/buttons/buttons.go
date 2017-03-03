package buttons

import (
	"./driver"
)

func Check_button_internal(button int) bool {
	press := driver.Elev_get_button_signal(int, int)
}

func Check_button_external(button int) bool {
	press := driver.Elev_get_button_signal(int, int)
}
