package main

import (
	"./definitions"
	"./master"
	"fmt"
)

func main() {
	elevatorState := definitions.ElevatorState{LastFloor: 8, Direction: -1, Destination: 4}
	internalPressOrder := definitions.Order{Floor: 10, Direction: 0}
	fmt.Println("elevatorState: ", elevatorState)

	fmt.Println("InternalPressorder:", internalPressOrder)
	orderList := definitions.Orders{Orders: []definitions.Order{definitions.Order{Floor: 9, Direction: -1}, definitions.Order{Floor: 4, Direction: -1}, definitions.Order{Floor: 5, Direction: -1}, definitions.Order{Floor: 9, Direction: -1}}}
	fmt.Println("OrderLIst: ", orderList)
	orderList = distributeInternalOrderToOrderList(internalPressOrder, orderList, elevatorState)
	fmt.Println("newOrderList: ", orderList)

	return
}

func distributeInternalOrderToOrderList(internalPressOrder definitions.Order, currentOrderList definitions.Orders, elevatorState definitions.ElevatorState) definitions.Orders {
	newOrderList := definitions.Orders{}

	if master.CheckForDuplicateOrder(currentOrderList, internalPressOrder.Floor) {
		return currentOrderList
	}

	tempNum := 0

	if elevatorState.Direction == 1 {
		// You are going up
		fmt.Println("You are going up")
		if currentOrderList.Orders[0].Floor == elevatorState.Destination { /* You can add in front of currentOrderList */
			fmt.Println("First order is the destination floor")
			newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
			copy(newOrderList.Orders[1:], newOrderList.Orders[:])
			newOrderList.Orders[0] = internalPressOrder
			return newOrderList
		} else { /* There are orders before destinationOrder */
			for i, order := range currentOrderList.Orders {
				fmt.Println("Order[", i, "]: ", order)
				if order.Floor > tempNum { // To check where you turn
					fmt.Println(" order.Floor > tempNum ")
					if order.Floor > internalPressOrder.Floor && elevatorState.LastFloor < internalPressOrder.Floor {
						fmt.Println("This IF STATEMENT")
						newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
						copy(newOrderList.Orders[i+1:], newOrderList.Orders[i:])
						newOrderList.Orders[i] = internalPressOrder
						return newOrderList
					}
					tempNum = order.Floor
				}
				if tempNum == elevatorState.Destination {
					fmt.Println("tempNum == elevatorState.Destination")

					for j, order2 := range currentOrderList.Orders {
						fmt.Println("Length, ", len(currentOrderList.Orders), ", j, ", j)
						if j > i {
							if order2.Floor < internalPressOrder.Floor {
								fmt.Println("The other IF STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+1:], newOrderList.Orders[j:])
								newOrderList.Orders[j] = internalPressOrder
								return newOrderList
							} else if j == len(currentOrderList.Orders)-1 {
								fmt.Println("This third STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+2:], newOrderList.Orders[j+1:])
								newOrderList.Orders[j+1] = internalPressOrder
								return newOrderList
							}
						}
					}
				}
			}
		}
	} else {
		tempNum = definitions.N_FLOORS
		fmt.Println("You are going down")
		if currentOrderList.Orders[0].Floor == elevatorState.Destination { /* You can add in front of currentOrderList */
			fmt.Println("First order is the destination floor")
			newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
			copy(newOrderList.Orders[1:], newOrderList.Orders[:])
			newOrderList.Orders[0] = internalPressOrder
			return newOrderList
		} else { /* There are orders before destinationOrder */
			for i, order := range currentOrderList.Orders {
				fmt.Println("Order[", i, "]: ", order)
				if order.Floor < tempNum { // To check where you turn
					fmt.Println(" order.Floor > tempNum ")
					if order.Floor < internalPressOrder.Floor && elevatorState.LastFloor < internalPressOrder.Floor {
						fmt.Println("This IF STATEMENT")
						newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
						copy(newOrderList.Orders[i+1:], newOrderList.Orders[i:])
						newOrderList.Orders[i] = internalPressOrder
						return newOrderList
					}
					tempNum = order.Floor
				}
				if tempNum == elevatorState.Destination {
					fmt.Println("tempNum == elevatorState.Destination")

					for j, order2 := range currentOrderList.Orders {
						fmt.Println("Length, ", len(currentOrderList.Orders), ", j, ", j)
						if j > i {
							if order2.Floor > internalPressOrder.Floor {
								fmt.Println("The other IF STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+1:], newOrderList.Orders[j:])
								newOrderList.Orders[j] = internalPressOrder
								return newOrderList
							} else if j == len(currentOrderList.Orders)-1 {
								fmt.Println("This third STATEMENT")

								newOrderList.Orders = append(currentOrderList.Orders, definitions.Order{})
								copy(newOrderList.Orders[j+2:], newOrderList.Orders[j+1:])
								newOrderList.Orders[j+1] = internalPressOrder
								return newOrderList
							}
						}
					}
				}
			}
		}
	}
	return definitions.Orders{}
}
