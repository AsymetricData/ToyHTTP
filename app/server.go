package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/routes"
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
	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error while reading Conn ", err)
	}
	//fmt.Println("Handled new data : ", n)

	router := routes.NewRouter("/", conn)
	router.Handle("/", func(conn net.Conn, r *request.Request) {
		writeResponse("HTTP/1.1 200 OK", 200, conn)
	})
	router.Handle("/echo/Coo/{value}", func(conn net.Conn, r *request.Request) {
		value := r.Params["value"]
		writeResponse("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: "+strconv.Itoa(len(value))+"\r\n\r\n"+value, 200, conn)
	})
	router.Handle("/echo/{value}", func(conn net.Conn, r *request.Request) {
		value := r.Params["value"]
		writeResponse("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: "+strconv.Itoa(len(value))+"\r\n\r\n"+value, 200, conn)
	})

	r := request.NewRequest(buffer)

	err = router.Get(&r)

	if err != nil {
		fmt.Println(err)
		writeResponse("HTTP/1.1 404 Not Found", 404, conn)
	}

}

func writeResponse(response string, status int, conn net.Conn) {

	//responseLen := len(response)

	_, err := conn.Write([]byte(response + "\r\n\r\n"))

	if err != nil {
		fmt.Println("Error Write ", err)
	}
}
