package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {
	// Define shell flags
	keyPath := flag.String("key", "", "Path to the SSH private key file (optional)")
	username := flag.String("user", "testuser", "Username for SSH login (default: testuser)")
	host := flag.String("host", "localhost", "Host and port to connect to (default: localhost:22)")
	flag.Parse()

	// Ensure the host has a port; if missing, append :22
	if !strings.Contains(*host, ":") {
		*host = *host + ":22"
	}

	// Define SSH authentication method
	var authMethod ssh.AuthMethod
	if *keyPath != "" {
		// Read the private key if a key path is specified
		key, err := os.ReadFile(*keyPath)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}

		// Parse the private key
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		// Use password if key path is not provided
		authMethod = ssh.Password("password123")
	}

	// Configure SSH client
	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{authMethod},
		// Use InsecureIgnoreHostKey for testing purposes; use FixedHostKey in production
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", *host, config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected to SSH. Enter commands to execute (type 'exit' to quit):")

	// Create a loop to read and execute commands interactively
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		command := scanner.Text()
		if strings.TrimSpace(command) == "exit" {
			fmt.Println("Exiting SSH session.")
			break
		}
		if command == "" {
			continue
		}

		// Create a new session for each command
		session, err := client.NewSession()
		if err != nil {
			log.Fatalf("Failed to create session: %v", err)
		}

		// Capture the command output
		output, err := session.CombinedOutput(command)
		if err != nil {
			log.Printf("Command execution error: %v\n", err)
		}

		fmt.Print(string(output))
		session.Close()
	}
}
