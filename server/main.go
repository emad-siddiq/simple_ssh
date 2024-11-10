// server/main.go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

func main() {
	// Load server private key
	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if string(pass) == "password123" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}
	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "0.0.0.0:2222")
	if err != nil {
		log.Fatalf("Failed to listen on 2222: %v", err)
	}

	log.Printf("Listening on 2222...")
	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			continue
		}
		go handleConnection(nConn, config)
	}
}

func handleConnection(nConn net.Conn, config *ssh.ServerConfig) {
	defer nConn.Close()

	conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("New SSH connection from %s (%s)", conn.RemoteAddr(), conn.ClientVersion())

	// Discard all global requests
	go ssh.DiscardRequests(reqs)

	// Accept all channels
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		// Handle session requests
		go handleChannelRequests(channel, requests)
	}
}

func executeCommand(command string, channel ssh.Channel) error {
	// Split the command and arguments
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	// Create the command
	cmd := exec.Command(parts[0], parts[1:]...)

	// Get pipes for stdin, stdout, and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Copy stdout to channel
	go func() {
		defer wg.Done()
		io.Copy(channel, stdout)
	}()

	// Copy stderr to channel
	go func() {
		defer wg.Done()
		io.Copy(channel, stderr)
	}()

	// Close stdin since we don't use it
	stdin.Close()

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		return err
	}

	// Wait for output copying to finish
	wg.Wait()

	return nil
}

func handleChannelRequests(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	var wg sync.WaitGroup

	for req := range requests {
		switch req.Type {
		case "pty-req":
			req.Reply(true, nil)
			log.Printf("PTY requested")

		case "shell":
			wg.Add(1)
			req.Reply(true, nil)
			log.Printf("Shell requested")

			go func() {
				defer wg.Done()
				fmt.Fprintf(channel, "\r\nWelcome to the SSH server!\r\n")
				fmt.Fprintf(channel, "Type 'exit' to close the session\r\n")

				for {
					fmt.Fprintf(channel, "\r\n$ ")
					buf := make([]byte, 1024)
					n, err := channel.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Printf("Error reading from channel: %v", err)
						}
						return
					}

					command := string(buf[:n])
					command = strings.TrimSpace(command)

					if command == "exit" {
						fmt.Fprintf(channel, "\r\nGoodbye!\r\n")
						return
					}

					// Execute the command and handle any errors
					if err := executeCommand(command, channel); err != nil {
						fmt.Fprintf(channel, "\r\nError executing command: %v\r\n", err)
					}
				}
			}()

		case "exec":
			var payload = struct{ Command string }{}
			ssh.Unmarshal(req.Payload, &payload)
			req.Reply(true, nil)
			executeCommand(payload.Command, channel)
			channel.Close()

		default:
			log.Printf("Received request type %v", req.Type)
			req.Reply(false, nil)
		}
	}

	wg.Wait()
}
