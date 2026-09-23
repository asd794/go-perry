package main

import (
	"bytes"
	"crypto/mlkem"
	"fmt"
	"log"
)

func main() {
	// === Scenario: Alice wants to securely share a key with Bob. ===

	// --- Alice's steps ---
	// 1. Alice generates an ML-KEM-768 private key (DecapsulationKey).
	//    GenerateKey768 returns (*DecapsulationKey768, error).
	privateKeyAlice, err := mlkem.GenerateKey768()
	if err != nil {
		log.Fatalf("Alice: Failed to generate ML-KEM-768 decapsulation key: %v", err)
	}

	// 2. Retrieve the corresponding public key (EncapsulationKey) from the private key.
	publicKeyAlice := privateKeyAlice.EncapsulationKey()

	// 3. Alice serializes her public key as a byte sequence to send to Bob.
	publicKeyAliceBytes := publicKeyAlice.Bytes()
	fmt.Printf("Alice's Public Key (ML-KEM-768, %d bytes): %x...\n", len(publicKeyAliceBytes), publicKeyAliceBytes[:16])

	// --- Bob's steps ---
	// Bob receives Alice's public key byte sequence, publicKeyAliceBytes.
	// 4. Bob creates a public key instance from the received byte sequence.
	publicKeyReceivedByBob, err := mlkem.NewEncapsulationKey768(publicKeyAliceBytes)
	if err != nil {
		log.Fatalf("Bob: Failed to parse Alice's public key: %v", err)
	}

	// 5. Bob uses Alice's public key to encapsulate a new shared key.
	sharedKeyForBob, ciphertextForAlice := publicKeyReceivedByBob.Encapsulate()
	fmt.Printf("Bob: Generated Shared Key (ML-KEM-768, %d bytes): %x\n", len(sharedKeyForBob), sharedKeyForBob)
	fmt.Printf("Bob: Generated Ciphertext for Alice (%d bytes): %x...\n", len(ciphertextForAlice), ciphertextForAlice[:16])

	// --- Alice's steps ---
	// Alice receives the ciphertext ciphertextForAlice from Bob.
	// 6. Alice uses her private key and the received ciphertext to decapsulate the shared key.
	sharedKeyForAlice, err := privateKeyAlice.Decapsulate(ciphertextForAlice)
	if err != nil {
		// Decapsulate returns an error if the ciphertext is invalid or has been tampered with.
		log.Fatalf("Alice: Failed to decapsulate shared key: %v", err)
	}
	fmt.Printf("Alice: Decapsulated Shared Key (ML-KEM-768, %d bytes): %x\n", len(sharedKeyForAlice), sharedKeyForAlice)

	// --- Verification ---
	// 7. Verify that Alice and Bob obtained the same shared key.
	if bytes.Equal(sharedKeyForAlice, sharedKeyForBob) {
		fmt.Println("\nSuccess! Alice and Bob now share the same secret key using ML-KEM-768.")
	} else {
		// This should not normally happen if decapsulation succeeds and the data has not been tampered with.
		fmt.Println("\nError! Shared keys do NOT match. This is unexpected.")
	}

}
