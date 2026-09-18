package main

/*
#cgo LDFLAGS: -luring
#include "uring.h"
#include <unistd.h>
*/
import "C"

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

func main() {
	serverFD := C.server_socket(8080)
	if serverFD < 0 {
		panic("failed to create server socket")
	}
	defer C.close(serverFD)

	ret := C.uring_init(256)
	if ret < 0 {
		panic("failed to initialize io_uring")
	}
	defer C.uring_cleanup()

	/*
		submit the first accept operation
	*/

	if ret := C.submit_accept(serverFD); ret < 0 {
		panic("failed to submit accept operation")
	}

	/*
		SIGINT
	*/

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("io_uring TCP echo server is running on port 8080. Press Ctrl+C to stop.")

	for {
		select {
		case <-sig:
			fmt.Println("Shutting down...")
			return
		default:
		}

		var completion C.completion_t

		ret := C.wait_cqe(&completion)
		if ret < 0 {
			fmt.Fprintf(os.Stderr, "wait_cqe failed: %v\n", ret)
			return
		}

		switch completion.op {
		case C.OP_ACCEPT:
			clientFD := completion.result
			if clientFD < 0 {
				fmt.Fprintf(os.Stderr, "accept failed: %v\n", clientFD)
				C.submit_accept(serverFD)
				continue
			}
			fmt.Printf("Accepted connection: fd=%d\n", clientFD)

			conn := C.connection_new(clientFD)

			if conn == nil {
				C.close(clientFD)
				C.submit_accept(serverFD)
				continue
			}

			if ret := C.submit_recv(conn); ret < 0 {
				C.connection_free(conn)
			}

			if ret := C.submit_accept(serverFD); ret < 0 {
				fmt.Fprintf(os.Stderr, "submit_accept failed: %v\n", ret)
			}

		case C.OP_RECV:
			conn := completion.conn
			n := completion.result

			if conn == nil {
				continue
			}

			if n == 0 {
				// fmt.Printf("Connection closed: fd=%d\n", conn.fd)
				C.connection_free(conn)
				continue
			}

			if n < 0 {
				fmt.Fprintf(os.Stderr, "recv failed: %v\n", n)
				C.connection_free(conn)
				continue
			}

			C.connection_set_len(conn, n)

			buf := C.GoBytes(unsafe.Pointer(C.connection_buffer(conn)), C.int(n))

			fmt.Printf("Received %d bytes: %s\n", n, string(buf))

			// fmt.Printf("Received %d bytes from fd=%d: %s\n", n, conn.fd, string(buf))

			if ret := C.submit_send(conn); ret < 0 {
				fmt.Fprintf(os.Stderr, "submit_send failed: %v\n", ret)
				C.connection_free(conn)
			}

		case C.OP_SEND:
			conn := completion.conn
			n := completion.result

			if n < 0 {
				fmt.Printf("send failed: %d\n", n)
				C.connection_free(conn)
				continue
			}

			C.connection_consume(conn, C.int(n))

			remaining := C.connection_remaining(conn)

			if remaining > 0 {
				ret := C.submit_send(conn)
				if ret < 0 {
					fmt.Printf("submit_send failed: %d\n", ret)
					C.connection_free(conn)
				}
			} else {
				ret := C.submit_recv(conn)
				if ret < 0 {
					fmt.Printf("submit_recv failed: %d\n", ret)
					C.connection_free(conn)
				}
			}
		}
	}
}
