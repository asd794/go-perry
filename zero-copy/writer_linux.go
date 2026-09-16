package main

import (
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

	file, err := os.OpenFile("/tmp/my_shared_mem", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	defer os.Remove("/tmp/my_shared_mem")
	if err := file.Truncate(int64(memSize)); err != nil {
		panic(err)
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, memSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)

	sharedPacket := (*Packet)(unsafe.Pointer(&data[0]))

	// Write the packet to shared memory
	sharedPacket.ID = 1
	sharedPacket.Score = 100
	sharedPacket.Ready = true

	time.Sleep(10 * time.Second)

}
