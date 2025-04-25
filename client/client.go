package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to Server.Type messages and press Enter to send.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Closing Connection.")
			return
		}
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing message:", err)
		}
	}

}
