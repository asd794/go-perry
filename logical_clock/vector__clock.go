package main

import "fmt"

type VectorClock struct {
	id    int
	clock []uint64
}

func NewVectorClock(id, n int) *VectorClock {
	return &VectorClock{
		id:    id,
		clock: make([]uint64, n),
	}
}

// Local event
func (v *VectorClock) Tick() {
	v.clock[v.id]++
}

// Send
func (v *VectorClock) Send() []uint64 {
	v.Tick()

	msg := make([]uint64, len(v.clock))
	copy(msg, v.clock)

	return msg
}

// Receive
func (v *VectorClock) Receive(msg []uint64) {
	for i := range v.clock {
		if msg[i] > v.clock[i] {
			v.clock[i] = msg[i]
		}
	}

	v.Tick()
}

func (v *VectorClock) String() string {
	return fmt.Sprint(v.clock)
}

func main() {
	A := NewVectorClock(0, 2)
	B := NewVectorClock(1, 2)

	// A local event
	A.Tick()

	fmt.Println("A:", A)
	// [1 0]

	// A -> B
	msg := A.Send()

	fmt.Println("A send:", msg)
	// [2 0]

	// B receives
	B.Receive(msg)

	fmt.Println("B:", B)
	// [2 1]

	// B local event
	B.Tick()

	fmt.Println("B:", B)
	// [2 2]
}
