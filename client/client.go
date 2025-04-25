package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Read server messages concurrently
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	// Send name
	serverReader := bufio.NewReader(conn)
	prompt, _ := serverReader.ReadString(':')
	fmt.Print(prompt)
	name, _ := reader.ReadString('\n')
	fmt.Fprint(conn, name)

	name = strings.TrimSpace(name)

	for {
		fmt.Println("\nChoose chat mode:\n1. Public\n2. Private\n3. Exit")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "3" {
			fmt.Fprintln(conn, "exit")
			break
		}

		switch choice {
		case "1":
			fmt.Print("Enter your message: ")
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			fmt.Fprintf(conn, "public|%s\n", msg)

		case "2":
			fmt.Print("Enter recipient name: ")
			recipient, _ := reader.ReadString('\n')
			recipient = strings.TrimSpace(recipient)
			fmt.Print("Enter your private message: ")
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			fmt.Fprintf(conn, "private|%s|%s\n", recipient, msg)

		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
