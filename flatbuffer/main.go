package main

import (
	"flatbuffer/generated/example"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
)

func main() {
	builder := flatbuffers.NewBuilder(0)

	name := builder.CreateString("Perry")

	example.UserStart(builder)
	example.UserAddId(builder, 1001)
	example.UserAddName(builder, name)
	example.UserAddAge(builder, 25)
	example.UserAddActive(builder, true)
	user := example.UserEnd(builder)

	builder.Finish(user)

	buf := builder.FinishedBytes()

	fmt.Printf("Binary size: %d bytes \n", len(buf))
	fmt.Printf("Binary data: %x \n\n", buf)

	userData := example.GetRootAsUser(buf, 0)

	fmt.Println("User:")
	fmt.Println("ID:", userData.Id())
	fmt.Println("Name:", string(userData.Name()))
	fmt.Println("Age:", userData.Age())
	fmt.Println("Active:", userData.Active())
}
