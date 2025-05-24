// Package main provides a utility to generate SHA-256 hashes for passwords
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

// HashPassword generates a SHA-256 hash for a password
func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	// Check if a password was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: hashpassword <password>")
		fmt.Println("Generates a SHA-256 hash for the provided password")
		os.Exit(1)
	}

	// Get the password from the command line arguments
	password := os.Args[1]

	// Hash the password
	hash := HashPassword(password)

	// Print the hash
	fmt.Printf("SHA-256 hash for '%s':\n%s\n", password, hash)
	fmt.Println("\nAdd this to your .env file as:")
	fmt.Printf("PRIVATE_USE_PASSWORD=%s\n", hash)
}
