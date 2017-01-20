package main

import (
	. "fmt"
	"runtime"
	"sync"
	//"time"
)

func thread1(i chan int, wg *sync.WaitGroup) {
	for j := 0; j < 1000000-1; j++ {
		local_i := <-i
		local_i++
		i <- local_i
	}
	wg.Done()
}

func thread2(i chan int, wg *sync.WaitGroup) {
	for j := 0; j < 1000000; j++ {
		local_i := <-i
		local_i--
		i <- local_i
	}
	wg.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	i := make(chan int, 2)
	i <- 0

	wg.Add(2)
	go thread1(i, &wg)
	go thread2(i, &wg)

	//time.Sleep(700 * time.Millisecond)
	wg.Wait()
	Println(<-i)
}
