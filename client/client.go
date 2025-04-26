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
