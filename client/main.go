package main

import (
	"flag"
	"fmt"
	"io"
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

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %v", err)
	}

	// Get stdin, stdout, and stderr pipes from the session
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("failed to get stdin pipe: %v", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("failed to get stdout pipe: %v", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatalf("failed to get stderr pipe: %v", err)
	}

	// Start the shell session
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %v", err)
	}

	// Print message to indicate shell is interactive
	fmt.Println("Shell started. You can now type commands...")

	// Connect input/output for interactive mode
	go handleShell(stdin, stdout, stderr)

	// Wait for session to finish (interactive shell will continue)
	if err := session.Wait(); err != nil {
		log.Fatalf("failed to wait for session: %v", err)
	}
}

// handleShell reads input from the terminal and sends it to the SSH session
func handleShell(stdin io.WriteCloser, stdout io.Reader, stderr io.Reader) {
	// Copy user input to stdin of the session
	go func() {
		_, err := io.Copy(stdin, os.Stdin)
		if err != nil {
			log.Fatalf("failed to copy input to stdin: %v", err)
		}
	}()

	// Copy stdout and stderr to the local terminal
	go func() {
		_, err := io.Copy(os.Stdout, stdout)
		if err != nil {
			log.Fatalf("failed to copy stdout: %v", err)
		}
	}()

	go func() {
		_, err := io.Copy(os.Stderr, stderr)
		if err != nil {
			log.Fatalf("failed to copy stderr: %v", err)
		}
	}()
}
