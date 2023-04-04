package main

import (
	"fmt"
	"time"
)

func countToTen(c chan int) {
	for i := range [10]int{} {
		time.Sleep(time.Second * 5)
		fmt.Printf("sending -> %d\n", i)
		c <- i
	}
}

func main() {
	c := make(chan int)
	go countToTen(c)

	for {
		a := <-c
		fmt.Printf("recieve -> %d\n", a)
		if c == nil {
			println("GO ROUTINE FINISH")
		}
	}

}
