// client.go
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Get the server address from the environment variable
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		fmt.Println("Error: SERVER_ADDRESS not set in .env file")
		return
	}

	// Setup TLS config (important for self-signed certificates)
	config := &tls.Config{
		InsecureSkipVerify: true, // ⚡ Important for self-signed certs (in production, use real validation)
	}

	// Dial the server using the address from the environment variable
	conn, err := tls.Dial("tcp", serverAddress, config)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	prompt, _ := serverReader.ReadString(':')
	fmt.Print(prompt)
	name, _ := reader.ReadString('\n')
	fmt.Fprint(conn, name)

	for {
		fmt.Println("\nChoose chat mode:")
		fmt.Println("1. Public")
		fmt.Println("2. Private")
		fmt.Println("3. Create Group")
		fmt.Println("4. Send Group Message")
		fmt.Println("5. Exit")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter your message: ")
			msg, _ := reader.ReadString('\n')
			fmt.Fprintf(conn, "public|%s\n", strings.TrimSpace(msg))

		case "2":
			fmt.Print("Enter recipient name: ")
			recipient, _ := reader.ReadString('\n')
			fmt.Print("Enter your private message: ")
			msg, _ := reader.ReadString('\n')
			fmt.Fprintf(conn, "private|%s|%s\n", strings.TrimSpace(recipient), strings.TrimSpace(msg))

		case "3":
			fmt.Print("Enter new group name: ")
			groupName, _ := reader.ReadString('\n')
			groupName = strings.TrimSpace(groupName)
			fmt.Fprintf(conn, "creategroup|%s\n", groupName)

			memberCount := 1
			for {
				fmt.Printf("Add member %d (or type 'done' to finish): ", memberCount)
				memberName, _ := reader.ReadString('\n')
				memberName = strings.TrimSpace(memberName)
				fmt.Fprintln(conn, memberName)

				if memberName == "done" {
					break
				}
				memberCount++
			}

		case "4":
			fmt.Print("Enter group name: ")
			groupName, _ := reader.ReadString('\n')
			fmt.Print("Enter group message: ")
			msg, _ := reader.ReadString('\n')
			fmt.Fprintf(conn, "groupmsg|%s|%s\n", strings.TrimSpace(groupName), strings.TrimSpace(msg))

		case "5":
			fmt.Fprintln(conn, "exit")
			return

		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
