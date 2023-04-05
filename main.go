package main

import (
	"fmt"
	"time"
)

func countToTen(c chan<- int) {
	for i := range [10]int{} {
		time.Sleep(time.Second * 2)
		fmt.Printf("sending -> %d\n", i)
		c <- i
	}
	close(c)
}

func main() {
	c := make(chan int)
	go countToTen(c)

	receive(c)

}

func receive(c <-chan int) {
	for {
		a, ok := <-c
		if ok != true {
			println("GO ROUTINE FINISH")
			break
		}
		fmt.Printf("recieve -> %d\n", a)

	}
}
