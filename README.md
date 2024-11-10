# Go SSH Client/Server Implementation

A simple SSH client and server implementation that supports connections to remote SSH servers like AWS or the bundled local server. It offers remote command execution and interactive shell sessions.

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
git clone https://github.com/emad-siddiq/simple_ssh
cd simple_ssh
```


1. Generate SSH key for the server:
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


## Security Considerations

This implementation is for educational purposes and includes several security compromises:

1. Uses hardcoded credentials
2. Uses `InsecureIgnoreHostKey` in the client
3. No encryption of sensitive configuration
4. No rate limiting
5. No audit logging
6. No command sanitization

For production use, one should:

- Implement proper user authentication
- Use proper host key verification
- Implement proper logging
- Add timeout mechanisms
- Add resource limits
- Implement proper error handling



## License

This project is licensed under the MIT License.
