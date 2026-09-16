package main

import (
	"errors"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	const memSize = 1024

	bufferName, _ := windows.UTF16PtrFromString("Local\\MySharedMemory")

	hMapFile, err := windows.CreateFileMapping(
		windows.InvalidHandle,  // use the system paging file
		nil,                    // default security attributes
		windows.PAGE_READWRITE, // read/write access
		0,                      // high-order DWORD of the maximum size of the shared memory segment
		memSize,                // size of the shared memory segment
		bufferName,             // shared memory name
	)
	if err != nil && !errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		panic(err)
	}

	defer windows.CloseHandle(hMapFile)

	addr, err := windows.MapViewOfFile(
		hMapFile,
		windows.FILE_MAP_WRITE,
		0,
		0,
		memSize,
	)
	if err != nil {
		panic(err)
	}
	defer windows.UnmapViewOfFile(addr)

	sharedSlice := (*[memSize]byte)(unsafe.Pointer(addr))

	message := []byte("Hello, shared memory!")
	copy(sharedSlice[:], message)

	time.Sleep(10 * time.Second)
}
