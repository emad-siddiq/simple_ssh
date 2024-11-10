# Go SSH Client/Server Implementation

A simple SSH client and server implementation in Go that allows for remote command execution and interactive shell sessions.

## Features

- Password-based authentication
- Interactive shell support
- Command execution
- PTY (pseudo-terminal) support
- Configurable port (default: 2222)
- Support for multiple concurrent connections

## Prerequisites

- Go 1.16 or later
- `golang.org/x/crypto/ssh` package

## Project Structure

```
go-ssh/
├── server/
│   ├── main.go
│   └── id_rsa     # Server private key (will be generated)
├── client/
│   └── main.go
└── README.md
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-ssh
```

2. Install dependencies:
```bash
go get golang.org/x/crypto/ssh
```

3. Generate SSH key for the server:
```bash
# From the project root directory
ssh-keygen -t rsa -f server/id_rsa
```

## Building

### Server
```bash
cd server
go build -o ssh-server
```

### Client
```bash
cd client
go build -o ssh-client
```

## Running

1. Start the server:
```bash
# If built:
./server/ssh-server

# Or run directly:
cd server
go run main.go
```

2. In another terminal, start the client:
```bash
# If built:
./client/ssh-client

# Or run directly:
cd client
go run main.go
```

## Default Configuration

- Server port: 2222
- Default username: "testuser"
- Default password: "password123"
- Server host: localhost

## Usage Examples

After connecting with the client, you can run various shell commands:

```bash
$ ls -la
$ pwd
$ date
$ echo "Hello, World!"
$ exit    # to close the session
```

## Modifying Default Settings

### Changing the Port

In `server/main.go`, modify the port in the `net.Listen` call:
```go
listener, err := net.Listen("tcp", "0.0.0.0:2222")  // Change 2222 to your desired port
```

In `client/main.go`, update the port in the `ssh.Dial` call:
```go
client, err := ssh.Dial("tcp", "localhost:2222", config)  // Change 2222 to match server port
```

### Changing Authentication

In `server/main.go`, modify the `PasswordCallback` function:
```go
PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
    if string(pass) == "your-new-password" {  // Change password here
        return nil, nil
    }
    return nil, fmt.Errorf("password rejected for %q", c.User())
},
```

In `client/main.go`, update the password in the `ClientConfig`:
```go
config := &ssh.ClientConfig{
    User: "your-username",  // Change username here
    Auth: []ssh.AuthMethod{
        ssh.Password("your-new-password"),  // Change password here
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
}
```

## Security Considerations

This implementation is for educational purposes and includes several security compromises:

1. Uses hardcoded credentials
2. Uses `InsecureIgnoreHostKey` in the client
3. No encryption of sensitive configuration
4. No rate limiting
5. No audit logging
6. No command sanitization

For production use, you should:

- Implement proper user authentication
- Use proper host key verification
- Add command whitelisting/blacklisting
- Implement proper logging
- Add timeout mechanisms
- Add resource limits
- Implement proper error handling

## Troubleshooting

1. **Port Already in Use**
```bash
Error: listen tcp :2222: bind: address already in use
```
Solution: Change the port number or stop the process using port 2222

2. **Permission Denied**
```bash
Error: permission denied
```
Solution: Check that the password in client matches the server configuration

3. **Private Key Not Found**
```bash
Error: Failed to load private key
```
Solution: Ensure you've generated the SSH key and placed it in the server directory


## License

This project is licensed under the MIT License.