package main

import (
	"fmt"
	"sync"
	// "time"
)

func thread1(c chan int, wg *sync.WaitGroup) {
	for j := 0; j < 1000000; j++ {
		i := <-c
		i++
		c <- i
	}
	wg.Done()
}

func thread2(c chan int, wg *sync.WaitGroup) {
	for j := 0; j < 1000000; j++ {
		i := <-c
		i--
		c <- i
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	fmt.Println("Hello Test")
	ch := make(chan int, 2)
	ch <- 0

	wg.Add(2)
	go thread1(ch, &wg)
	go thread2(ch, &wg)
	// time.Sleep(1400 * time.Millisecond)
	// x, y := <-ch, <-ch //Recieving from
	wg.Wait()
	fmt.Println(<-ch)
}
