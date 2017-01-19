package main

import (
	. "fmt"
	"runtime"
	// "time"
)

func someGoroutine() {
	Println("Hello from a goRoutine!")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Println("Hello World1!")
	go someGoroutine()
	// time.Sleep(2000 * time.Millisecond)
	Println("Hello World2!")

}

//COMMENT
