package main

import (
	. "fmt"
	//"runtime"
	"time"
)

func someGoroutine() {
	Println("Hello from a goroutine")
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	go someGoroutine()
	time.Sleep(100 * time.Millisecond)
}
