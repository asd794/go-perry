package main

import "fmt"

type LamportClock struct {
	time uint64
}

func (c *LamportClock) Tick() uint64 {
	c.time++
	return c.time
}

func (c *LamportClock) Send() uint64 {
	return c.Tick()
}

func (c *LamportClock) Receive(received uint64) uint64 {
	if received > c.time {
		c.time = received
	}

	c.time++
	return c.time
}

func main() {
	A := &LamportClock{}
	B := &LamportClock{}

	// A: event
	fmt.Println("A event:", A.Tick()) // 1

	// A -> B
	msg := A.Send()
	fmt.Println("A send:", msg) // 2

	// B receives
	fmt.Println("B receive:", B.Receive(msg)) // 3

	// B local event
	fmt.Println("B event:", B.Tick()) // 4

	// B -> A
	msg = B.Send()
	fmt.Println("B send:", msg) // 5

	// A receives
	fmt.Println("A receive:", A.Receive(msg)) // 6
}
