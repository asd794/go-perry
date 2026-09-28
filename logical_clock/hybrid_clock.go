package main

import (
	"fmt"
	"sync"
	"time"
)

type Timestamp struct {
	Wall    int64  // physical time
	Logical uint64 // logical time (while physical time is the same)
}

type HLC struct {
	wall    int64
	logical uint64
}

func NewHLC() *HLC {
	return &HLC{}
}

func (h *HLC) Now() Timestamp {
	now := time.Now().UnixNano()

	switch {
	case now > h.wall:
		h.wall = now
		h.logical = 0
	default:
		h.logical++
	}

	return Timestamp{
		Wall:    h.wall,
		Logical: h.logical,
	}
}

func (h *HLC) Update(remote Timestamp) Timestamp {
	now := time.Now().UnixNano()

	maxWall := max(now, h.wall, remote.Wall)

	switch {
	case maxWall == h.wall && maxWall == remote.Wall:
		h.logical = max(h.logical, remote.Logical) + 1
	case maxWall == h.wall:
		h.logical++
	case maxWall == remote.Wall:
		h.logical = remote.Logical + 1
	default:
		h.logical = 0
	}

	h.wall = maxWall

	return Timestamp{
		Wall:    h.wall,
		Logical: h.logical,
	}
}

func main() {
	nodeA := NewHLC()
	nodeB := NewHLC()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		ts := nodeA.Now()

		fmt.Println("A local event:", ts)

		time.Sleep(10 * time.Millisecond)

		ts = nodeA.Now()

		fmt.Println("A local event:", ts)
	}()

	go func() {
		defer wg.Done()

		ts := nodeB.Now()

		fmt.Println("B local event:", ts)
	}()

	wg.Wait()

	// Simulate:
	//
	// A sends a message to B.
	//

	a := nodeA.Now()

	fmt.Println("\nA sends message:")
	fmt.Println("A:", a)

	b := nodeB.Update(a)

	fmt.Println("B receives message:")
	fmt.Println("B:", b)

	// B creates another local event.
	b = nodeB.Now()

	fmt.Println("B local event:")
	fmt.Println("B:", b)
}
