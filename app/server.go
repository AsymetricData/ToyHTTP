package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind port 4221")
		os.Exit(1)
	}
	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Failed to accept ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error while reading Conn ", err)
		}

		handleRequest(conn, buffer, n)
	}
}

func handleRequest(conn net.Conn, buffer []byte, n int) {
	fmt.Println("Handled new data : ", n)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
