package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func main() {
	file, err := os.Open("data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)

	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		panic(err)
	}

	err = rawConn.Control(func(fd uintptr) {
		var offset int64

		_, err := syscall.Sendfile(
			int(fd),
			int(file.Fd()),
			&offset,
			1024*1024,
		)
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("file sent")
}
