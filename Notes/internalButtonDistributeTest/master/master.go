package master

import (
	"../definitions"
	"fmt"
)

// Don't accept more orders to same floor. Assume every person gets on elevator.
func CheckForDuplicateOrder(orders definitions.Orders, buttonPressedFloor int) bool {
	for i := range orders.Orders {
		if orders.Orders[i].Floor == buttonPressedFloor {
			return true
		}
	}
	return false
}

func FloorIsInbetween(orderFloor int, buttonFloor int, elevatorLastFloor int, elevatorDirection int) bool {
	switch elevatorDirection {
	case definitions.DIR_UP:
		return buttonFloor > elevatorLastFloor &&
			buttonFloor < orderFloor
	case definitions.DIR_DOWN:
		return buttonFloor < elevatorLastFloor &&
			buttonFloor > orderFloor
	default:
		fmt.Println("Something is wrong in floorIsBetween()")
		return false
	}
}
