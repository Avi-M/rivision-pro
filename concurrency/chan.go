package main

import "fmt"
import "runtime"

func main() {
	// Goroutine num includes main processing
	fmt.Println(runtime.NumGoroutine())

	// Spawn two goroutines
	go func() {}()
	go func() {}()

	// Total three goroutines run
	fmt.Println(runtime.NumGoroutine()) 
}