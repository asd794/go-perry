//go:build linux

package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

const (
	addr    = ":8080"
	file    = "data.txt"
	bufSize = 64 * 1024
)

func main() {
	src, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("splice server listening on", addr)

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

	var transferErr error

	err = rawConn.Control(func(socketFD uintptr) {
		var pipe [2]int

		if err := unix.Pipe(pipe[:]); err != nil {
			transferErr = err
			return
		}

		defer unix.Close(pipe[0])
		defer unix.Close(pipe[1])

		for {
			// File -> Kernel Pipe
			n, err := unix.Splice(
				int(src.Fd()),
				nil,
				pipe[1],
				nil,
				bufSize,
				0,
			)
			if err != nil {
				transferErr = fmt.Errorf("file -> pipe: %w", err)
				return
			}

			if n == 0 {
				break
			}

			remaining := n

			// Kernel Pipe -> TCP Socket
			for remaining > 0 {
				written, err := unix.Splice(
					pipe[0],
					nil,
					int(socketFD),
					nil,
					int(remaining),
					0,
				)
				if err != nil {
					transferErr = fmt.Errorf("pipe -> socket: %w", err)
					return
				}

				if written == 0 {
					transferErr = fmt.Errorf("splice wrote 0 bytes")
					return
				}

				remaining -= written
			}

			fmt.Printf("spliced %d bytes\n", n)
		}
	})

	if err != nil {
		panic(err)
	}

	if transferErr != nil {
		panic(transferErr)
	}

	fmt.Println("transfer completed")
}
