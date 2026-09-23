package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"crypto/mldsa"
)

func main() {
	// 1. Select the parameter set and generate an ML-DSA-44 post-quantum key pair.
	params := mldsa.MLDSA44()
	priv, err := mldsa.GenerateKey(params)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	// 2. Prepare the message to sign.
	msg := []byte("This is an important message that needs post-quantum protection.")

	// 3. Create a digital signature.
	// opt := &mldsa.Options{
	// 	Context: "document-signature",
	// }
	// Use the options for signing.
	sig, err := priv.Sign(rand.Reader, msg, &mldsa.Options{})
	if err != nil {
		log.Fatalf("Failed to sign message: %v", err)
	}
	fmt.Printf("Signature size: %d bytes\n", len(sig))
	fmt.Printf("First 32 bytes of signature: %x\n", sig[:32])

	// 4. Verify the signature.
	pub := priv.PublicKey()
	err = mldsa.Verify(pub, msg, sig, &mldsa.Options{})
	if err != nil {
		fmt.Println("Signature verification failed:", err)
	} else {
		fmt.Println("Signature verification succeeded!")
	}
}
