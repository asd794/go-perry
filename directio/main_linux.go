package main

import (
	"errors"
	"os"
	"syscall"
)

const (
	AigenSzie = 4096 // NVMe SSD block size
)

func MakeAlignedBuffer(size int) ([]byte, error) {
	if size%AigenSzie != 0 {
		return nil, errors.New("size must be a multiple of AigenSzie")
	}
	block, err := syscall.Mmap(-1, 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANONYMOUS|syscall.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func FreeAlignedBuffer(b []byte) error {
	return syscall.Munmap(b)
}

func OpenDirectFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag|syscall.O_DIRECT, perm)
}
