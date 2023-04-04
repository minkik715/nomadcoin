package main

import (
	"fmt"
	"time"
)

func countToTen(name string) {
	for i := range [10]int{} {
		fmt.Printf("%s, -> %d \n", name, i)
		time.Sleep(time.Second * 1)
	}
}

func main() {
	go countToTen("FIRST")
	go countToTen("SECOND")
	for {

	}
}
