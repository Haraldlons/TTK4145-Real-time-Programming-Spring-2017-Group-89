package main

import "fmt"

func main() {
	ch := make(chan int, 9)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	ch <- 6
	ch <- 7
	ch <- 8
	for i := 0; i < 8; i++ {

		fmt.Println(<-ch)
	}
}
