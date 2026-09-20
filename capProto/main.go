package main

import (
	generated "capproto/schema"
	"fmt"
	"log"

	"capnproto.org/go/capnp/v3"
)

func main() {
	// Create a new Cap'n Proto message.
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		log.Fatal(err)
	}

	// Create the root User object.
	user, err := generated.NewRootUser(seg)
	if err != nil {
		log.Fatal(err)
	}

	// Set fields.
	user.SetId(1001)
	user.SetName("Perry")
	user.SetAge(25)
	user.SetActive(true)

	// Serialize the message.
	buf, err := msg.Marshal()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Binary size: %d bytes\n", len(buf))
	fmt.Printf("Binary data: %x\n\n", buf)

	// Read the serialized message.
	msg2, err := capnp.Unmarshal(buf)
	if err != nil {
		log.Fatal(err)
	}

	user2, err := generated.ReadRootUser(msg2)
	if err != nil {
		log.Fatal(err)
	}

	userName, _ := user2.Name()

	// Access fields directly through the generated API.
	fmt.Println("User:")
	fmt.Println("ID:", user2.Id())
	fmt.Println("Name:", userName)
	fmt.Println("Age:", user2.Age())
	fmt.Println("Active:", user2.Active())
}
