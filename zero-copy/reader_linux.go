package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type Packet struct {
	ID    uint32
	Score uint32
	Ready bool
}

func main() {
	var p Packet

	memSize := int(unsafe.Sizeof(p))

	file, err := os.OpenFile("/tmp/my_shared_mem", os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := syscall.Mmap(int(file.Fd()), 0, memSize, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)

	sharedPacket := (*Packet)(unsafe.Pointer(&data[0]))

	for !sharedPacket.Ready {
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("ID: %d, Score: %d\n", sharedPacket.ID, sharedPacket.Score)

}
