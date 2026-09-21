package main

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/fxamacker/cbor/v2"
)

type User struct {
	ID    uint64 `cbor:"id"`
	Name  string `cbor:"name"`
	Admin bool   `cbor:"admin"`
}

// func main() {
// 	user := User{
// 		ID:    1001,
// 		Name:  "Perry",
// 		Admin: true,
// 	}

// 	data, err := cbor.Marshal(user)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("CBOR bytes: %x\n", data)

// 	var decoded User

// 	if err := cbor.Unmarshal(data, &decoded); err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("Decoded: %+v\n", decoded)
// }

// Canonical

func main() {
	valueA := map[string]any{
		"name":   "Perry",
		"age":    uint64(25),
		"active": true,
	}

	valueB := map[string]any{
		"active": true,
		"name":   "Perry",
		"age":    uint64(25),
	}

	encMode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		log.Fatal(err)
	}

	encodedA, err := encMode.Marshal(valueA)
	if err != nil {
		log.Fatal(err)
	}

	encodedB, err := encMode.Marshal(valueB)
	if err != nil {
		log.Fatal(err)
	}

	hashA := sha256.Sum256(encodedA)
	hashB := sha256.Sum256(encodedB)

	fmt.Printf("Canonical CBOR A: %x\n", encodedA)
	fmt.Printf("SHA-256        A:        %x\n", hashA)

	fmt.Printf("Canonical CBOR B: %x\n", encodedB)
	fmt.Printf("SHA-256        B:        %x\n", hashB)
}
