package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {

	data, err := os.ReadFile("test.crt")
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		panic("failed to parse PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	fmt.Println("Issuer:", cert.Issuer)
	fmt.Println("Subject:", cert.Subject)
	fmt.Println("Common Name:", cert.Subject.CommonName)
	fmt.Println("Not Before:", cert.NotBefore)
	fmt.Println("Not After:", cert.NotAfter)
	fmt.Println("Serial Number:", cert.SerialNumber)
}
