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

	path := getPath(buffer)

	switch path {
	case "/":
		writeResponse("HTTP/1.1 200 OK", conn)
	default:
		writeResponse("HTTP/1.1 404 Not Found", conn)

	}
}

func writeResponse(response string, conn net.Conn) {
	_, err := conn.Write([]byte(response + "\r\n\r\n"))

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
				//HTTPmethod = string(slices[0])
				HTTPpath = string(slices[1])
				return HTTPpath
			}
		}

		sn = lineIndex + 2
	}

	return HTTPpath
}
