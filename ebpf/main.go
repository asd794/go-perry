package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
)

//go:generate bpf2go -cc clang -cflags "-I/usr/include/x86_64-linux-gnu" bpf bpf/tracker.c
//go:generate bpf2go -cc clang -cflags "-I/usr/include/x86_64-linux-gnu" packet bpf/packet.c

func main() {
	// Load process tracing eBPF program.
	// var trackerObjs bpfObjects
	// if err := loadBpfObjects(&trackerObjs, nil); err != nil {
	// 	log.Fatalf("loading tracker objects: %v", err)
	// }
	// defer trackerObjs.Close()

	// // Attach process tracing program to execve tracepoint.
	// tp, err := link.Tracepoint(
	// 	"syscalls",
	// 	"sys_enter_execve",
	// 	trackerObjs.TraceExecve,
	// 	nil,
	// )
	// if err != nil {
	// 	log.Fatalf("attaching execve tracepoint: %v", err)
	// }
	// defer tp.Close()

	// Get network interface.
	iface, err := net.InterfaceByName("wlan0")
	if err != nil {
		log.Fatalf("getting network interface: %v", err)
	}

	// Load packet/XDP eBPF program.
	var packetObjs packetObjects
	if err := loadPacketObjects(&packetObjs, nil); err != nil {
		log.Fatalf("loading packet objects: %v", err)
	}
	defer packetObjs.Close()

	// Attach XDP program to network interface.
	xdpLink, err := link.AttachXDP(link.XDPOptions{
		Program:   packetObjs.PacketFilter,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("attaching XDP: %v", err)
	}
	defer xdpLink.Close()

	fmt.Printf("eBPF programs loaded successfully\n")
	fmt.Printf("process tracing: syscalls/sys_enter_execve\n")
	fmt.Printf("packet filtering: XDP on %s\n", iface.Name)

	// Wait for Ctrl+C / SIGTERM.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	fmt.Println("exiting...")
}
