package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	const memSize = 1024

	bufferName, _ := windows.UTF16PtrFromString("Local\\MySharedMemory")

	hMapFile, _ := windows.CreateFileMapping(
		windows.InvalidHandle, // use the system paging file
		nil,                   // default security attributes
		windows.PAGE_READONLY, // read-only access
		0,                     // high-order DWORD of the maximum size of the shared memory segment
		memSize,               // size of the shared memory segment
		bufferName,            // shared memory name
	)

	defer windows.CloseHandle(hMapFile)

	addr, err := windows.MapViewOfFile(
		hMapFile,
		windows.FILE_MAP_READ,
		0,
		0,
		memSize,
	)
	if err != nil {
		panic(err)
	}
	defer windows.UnmapViewOfFile(addr)

	sharedSlice := (*[memSize]byte)(unsafe.Pointer(addr))
	println(string(sharedSlice[:]))
}
