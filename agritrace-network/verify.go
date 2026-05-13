package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: verify <ca_cert> <admin_cert>")
		os.Exit(1)
	}

	caCertPath := os.Args[1]
	adminCertPath := os.Args[2]

	caBytes, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Printf("Failed to read CA cert: %v\n", err)
		os.Exit(1)
	}

	adminBytes, err := ioutil.ReadFile(adminCertPath)
	if err != nil {
		fmt.Printf("Failed to read Admin cert: %v\n", err)
		os.Exit(1)
	}

	caBlock, _ := pem.Decode(caBytes)
	if caBlock == nil {
		fmt.Println("Failed to parse CA PEM")
		os.Exit(1)
	}

	adminBlock, _ := pem.Decode(adminBytes)
	if adminBlock == nil {
		fmt.Println("Failed to parse Admin PEM")
		os.Exit(1)
	}

	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		fmt.Printf("Failed to parse CA cert: %v\n", err)
		os.Exit(1)
	}

	adminCert, err := x509.ParseCertificate(adminBlock.Bytes)
	if err != nil {
		fmt.Printf("Failed to parse Admin cert: %v\n", err)
		os.Exit(1)
	}

	roots := x509.NewCertPool()
	roots.AddCert(caCert)

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	if _, err := adminCert.Verify(opts); err != nil {
		fmt.Printf("Verification failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verification successful!")
}
