# Multi-User Chat Application

A secure and scalable multi-user chat application built using Go, featuring **private**, **public**, and **group messaging** functionalities. The app uses **TLS encryption** to secure communication between clients and server.


## Features

- **Private Messaging**: Send direct messages to a specific user.
- **Public Messaging**: Broadcast messages to all connected users.
- **Group Messaging**: Create groups and chat within them.
- **TLS Encryption**: Secure all communication using SSL/TLS certificates and keys generated via OpenSSL.



## Tech Stack

- **Language**: Go (Golang)
- **Security**: TLS/SSL (OpenSSL Certificates)
- **Networking**: TCP Protocol



## Installation and Setup

### 1. Install Go

- Download and install Go from [https://golang.org/dl/](https://golang.org/dl/).
- Verify the installation:
  ```bash
  go version
  ```

### 2. Install OpenSSL and Generate Certificates

- **Install OpenSSL**:
  - **Windows**: [Download OpenSSL](https://slproweb.com/products/Win32OpenSSL.html)
  - **Linux (Ubuntu)**:
    ```bash
    sudo apt update
    sudo apt install openssl
    ```
  - **MacOS**:
    ```bash
    brew install openssl
    ```

- **Generate TLS Certificate and Private Key**:
  ```bash
  openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
  ```
  - `server.key` → Private Key
  - `server.crt` → Public Certificate

- Place `server.crt` and `server.key` inside your project folder.

### 3. Clone the Repository

```bash
git clone https://github.com/msaakaash/multiuser-chat-app.git
cd multiuser-chat-app
```

### 4. Running the Server

```bash
go run server.go
```

### 5. Running the Client(s)

```bash
go run client.go
```

You can run multiple clients to simulate multiple users chatting.

## Code of Conduct

Please read our [Code of Conduct](./CODE_OF_CONDUCT.md) before contributing to this project.


## Contributing

We welcome contributions from everyone!

1. **Fork** this repository.
2. **Create** a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Commit** your changes:
   ```bash
   git commit -m "Add your message"
   ```
4. **Push** the changes:
   ```bash
   git push origin feature/your-feature-name
   ```
5. **Open** a Pull Request.


## License  
This project is licensed under the [MIT License](LICENSE).



## Additional Notes

- Ensure the server is running before clients attempt to connect.
- The communication between client and server is fully encrypted with TLS.
- Groups can be created dynamically and members can join iteratively.


 ## Author
 [**Aakaash M S**](https://github.com/msaakaash)

