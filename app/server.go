package main

import (
	"bytes"
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
	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Failed to accept ", err.Error())
			os.Exit(1)
		}

		handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error while reading Conn ", err)
	}
	fmt.Println("Handled new data : ", n)

	getPath(buffer)

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	if err != nil {
		fmt.Println("Error Write ", err)
	}
}

func getPath(buffer []byte) string {
	//HTTPmethod := ""
	HTTPpath := ""
	sn := 0
	//protocol := ""

	for sn <= len(buffer) {
		lineIndex := bytes.Index(buffer, []byte("\r\n"))

		if lineIndex == -1 {
			fmt.Println("Error in the header")
			os.Exit(1)
		}

		line := buffer[sn:lineIndex]

		for lineIndex != -1 {
			slices := bytes.Split(line, []byte(" "))

			if len(slices) != 3 {
				break
			}

			if string(slices[0]) == "GET" || string(slices[0]) == "POST" {
				fmt.Println("Ok")
			}
		}

		sn = lineIndex + 2
	}

	return HTTPpath
}
