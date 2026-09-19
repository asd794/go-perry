package main

import (
	"fmt"
	"log"

	"github.com/vmihailenco/msgpack/v5"
)

type User struct {
	ID    int    `msgpack:"id"`
	Name  string `msgpack:"name"`
	Email string `msgpack:"email"`
	Age   int    `msgpack:"age"`
}

func main() {
	user := User{
		ID:    1001,
		Name:  "Perry",
		Email: "perry@example.com",
		Age:   25,
	}

	data, err := msgpack.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MessagePack bytes: %x\n", data)
	fmt.Printf("Size: %d bytes\n", len(data))

	// Decode
	var decoded User

	err = msgpack.Unmarshal(data, &decoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded: %+v\n", decoded)
}
