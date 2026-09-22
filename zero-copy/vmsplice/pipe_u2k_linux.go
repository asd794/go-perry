//go:build linux

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func main() {
	var pipe [2]int

	if err := unix.Pipe(pipe[:]); err != nil {
		panic(err)
	}

	defer unix.Close(pipe[0])
	defer unix.Close(pipe[1])

	data := []byte("hello from vmsplice\n")

	iov := unix.Iovec{
		Base: &data[0],
		Len:  uint64(len(data)),
	}

	n, err := unix.Vmsplice(
		pipe[1],
		[]unix.Iovec{iov},
		0,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("vmsplice: %d bytes\n", n)

	buf := make([]byte, len(data))

	nRead, err := unix.Read(
		pipe[0],
		buf,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("read: %q\n", buf[:nRead])
}
