package main

import (
	. "fmt"
	"runtime"
	"time"
)

//* is a reference
func thread1(i *int) {
	for j := 0; j < 1000000; j++ {
		*i++
	}
}

func thread2(i *int) {
	for j := 0; j < 1000000; j++ {
		*i--
	}
}

func someGoroutine() {
	Println("Hello from a goRoutine!")
}

func main() {
	//
	i float64 = 0 
	runtime.GOMAXPROCS(runtime.NumCPU())
	Println("Hello World1!")
	go someGoroutine()
	time.Sleep(100 * time.Millisecond)
	go thread1(&i)
	go thread2(&i)
	time.Sleep(100 * time.Millisecond)
	// Println("Hello World2!")
	Println("i: ", i)

}
