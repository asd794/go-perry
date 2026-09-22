//go:build linux

package main

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

const (
	addr = "127.0.0.1:9090"
	size = 64 * 1024
)

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)

	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		panic(err)
	}

	data := make([]byte, size)

	for i := range data {
		data[i] = 'A'
	}

	var sendErr error

	err = rawConn.Control(func(fd uintptr) {
		socketFD := int(fd)

		// Enable SO_ZEROCOPY.
		sendErr = unix.SetsockoptInt(
			socketFD,
			unix.SOL_SOCKET,
			unix.SO_ZEROCOPY,
			1,
		)
		if sendErr != nil {
			return
		}

		// MSG_ZEROCOPY asks Linux to avoid copying
		// the userspace buffer into the socket buffer
		// when possible.
		_, sendErr = unix.SendmsgN(
			socketFD,
			data,
			nil,
			nil,
			unix.MSG_ZEROCOPY,
		)
	})

	if err != nil {
		panic(err)
	}

	if sendErr != nil {
		panic(sendErr)
	}

	fmt.Printf("sent %d bytes with MSG_ZEROCOPY\n", len(data))

	// Keep the process alive briefly so that the kernel
	// has time to process the transmission.
	time.Sleep(time.Second)
}
