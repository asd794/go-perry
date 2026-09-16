package main

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkBufferedWrite(b *testing.B) {
	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "buffered.bin")

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		b.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	data := make([]byte, AigenSzie)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Write(data)
		if err != nil {
			b.Fatalf("Failed to write to file: %v", err)
		}
	}
}

func BenchmarkDirectIOWrite(b *testing.B) {
	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "directio.bin")

	f, err := OpenDirectFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		b.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	data, err := MakeAlignedBuffer(AigenSzie)
	if err != nil {
		b.Fatalf("Failed to create aligned buffer: %v", err)
	}
	defer FreeAlignedBuffer(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Write(data)
		if err != nil {
			b.Fatalf("Failed to write to file: %v", err)
		}
	}
}
