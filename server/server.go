package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients     = make(map[string]net.Conn)
	clientsLock = sync.Mutex{}
)

func listMembers() string {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	names := []string{}
	for name := range clients {
		names = append(names, name)
	}
	return "Active members: " + strings.Join(names, ", ")
}

func broadcast(sender, message string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	for name, conn := range clients {
		if name != sender {
			fmt.Fprintf(conn, "%s: %s\n", sender, message)
		}
	}
}

func sendPrivate(sender, recipient, message string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	if conn, ok := clients[recipient]; ok {
		fmt.Fprintf(conn, "[PM from %s]: %s\n", sender, message)
	} else {
		if conn, ok := clients[sender]; ok {
			fmt.Fprintf(conn, "User %s not found\n", recipient)
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	conn.Write([]byte("Enter your name: "))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	clientsLock.Lock()
	if _, exists := clients[name]; exists {
		conn.Write([]byte("Name already taken. Disconnecting...\n"))
		clientsLock.Unlock()
		return
	}
	clients[name] = conn
	clientsLock.Unlock()

	fmt.Printf("%s has joined the chat\n", name)
	broadcast("Server", fmt.Sprintf("%s has joined the chat", name))
	broadcast("Server", listMembers())

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)

		if message == "exit" {
			break
		}

		// Format: public|message or private|recipient|message
		if strings.HasPrefix(message, "public|") {
			msg := strings.TrimPrefix(message, "public|")
			broadcast(name, msg)
		} else if strings.HasPrefix(message, "private|") {
			parts := strings.SplitN(strings.TrimPrefix(message, "private|"), "|", 2)
			if len(parts) == 2 {
				recipient := parts[0]
				msg := parts[1]
				sendPrivate(name, recipient, msg)
			} else {
				conn.Write([]byte("Invalid private message format.\n"))
			}
		} else {
			conn.Write([]byte("Invalid format. Use public|message or private|recipient|message\n"))
		}
	}

	clientsLock.Lock()
	delete(clients, name)
	clientsLock.Unlock()
	fmt.Printf("%s has left the chat\n", name)
	broadcast("Server", fmt.Sprintf("%s has left the chat", name))
	broadcast("Server", listMembers())
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
