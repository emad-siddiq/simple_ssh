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
./client/ssh-client -key ~/.ssh/aws.pem -user ubuntu -host <aws-host-url>

# Or run directly:
cd client
go run . -key ~/.ssh/aws.pem -user ubuntu -host <aws-host-url>
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



## License

This project is licensed under the MIT License.