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

	groups     = make(map[string][]string)
	groupsLock = sync.Mutex{}
)

func listMembers() string {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	var names []string
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
	} else if conn, ok := clients[sender]; ok {
		fmt.Fprintf(conn, "User %s not found\n", recipient)
	}
}

func sendGroupMessage(sender, groupName, message string) {
	groupsLock.Lock()
	members, exists := groups[groupName]
	groupsLock.Unlock()

	if !exists {
		if conn, ok := clients[sender]; ok {
			fmt.Fprintf(conn, "Group %s does not exist.\n", groupName)
		}
		return
	}

	for _, member := range members {
		if member != sender {
			if conn, ok := clients[member]; ok {
				fmt.Fprintf(conn, "[Group %s | %s]: %s\n", groupName, sender, message)
			}
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

		switch {
		case strings.HasPrefix(message, "public|"):
			msg := strings.TrimPrefix(message, "public|")
			broadcast(name, msg)

		case strings.HasPrefix(message, "private|"):
			parts := strings.SplitN(strings.TrimPrefix(message, "private|"), "|", 2)
			if len(parts) == 2 {
				sendPrivate(name, parts[0], parts[1])
			} else {
				conn.Write([]byte("Invalid private message format.\n"))
			}

		case strings.HasPrefix(message, "creategroup|"):
			groupName := strings.TrimPrefix(message, "creategroup|")

			groupsLock.Lock()
			groups[groupName] = []string{name}
			groupsLock.Unlock()
			conn.Write([]byte("Enter member names (type 'done' to finish):\n"))

			for {
				member, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				member = strings.TrimSpace(member)
				if member == "done" {
					break
				}

				clientsLock.Lock()
				_, exists := clients[member]
				clientsLock.Unlock()

				if exists {
					groupsLock.Lock()
					groups[groupName] = append(groups[groupName], member)
					groupsLock.Unlock()
					conn.Write([]byte(fmt.Sprintf("Added %s to group %s\n", member, groupName)))
				} else {
					conn.Write([]byte(fmt.Sprintf("User %s does not exist.\n", member)))
				}
			}

			conn.Write([]byte(fmt.Sprintf("Group %s created successfully.\n", groupName)))

		case strings.HasPrefix(message, "groupmsg|"):
			parts := strings.SplitN(strings.TrimPrefix(message, "groupmsg|"), "|", 2)
			if len(parts) == 2 {
				sendGroupMessage(name, parts[0], parts[1])
			} else {
				conn.Write([]byte("Invalid group message format.\n"))
			}

		default:
			conn.Write([]byte("Invalid format.\n"))
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
