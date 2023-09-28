package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/routes"
)

var staticDirectory string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Server is starting...")

	if len(os.Args) > 2 {
		if os.Args[1] == "--directory" {
			staticDirectory = os.Args[2]
		}
	}

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

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)

	//ok

	if err != nil {
		//fmt.Println("Error while reading Conn ", err)
		return
	}
	//fmt.Println("Handled new data : ", n)

	router := routes.NewRouter("/", conn)
	router.Handle("/", func(conn net.Conn, r *request.Request) {
		writeResponse("HTTP/1.1 200 OK", 200, conn)
	})
	router.Handle("/echo/{val}/{value}", func(conn net.Conn, r *request.Request) {
		value := r.Params["val"] + "/" + r.Params["value"]
		writeResponse("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: "+strconv.Itoa(len(value))+"\r\n\r\n"+value, 200, conn)
	})
	router.Handle("/echo/{value}", func(conn net.Conn, r *request.Request) {
		value := r.Params["value"]
		writeResponse("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: "+strconv.Itoa(len(value))+"\r\n\r\n"+value, 200, conn)
	})
	router.Handle("/user-agent", func(conn net.Conn, r *request.Request) {
		value := r.Headers.UserAgent
		writeResponse("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: "+strconv.Itoa(len(value))+"\r\n\r\n"+value, 200, conn)
	})
	router.Handle("/files/{value}", func(conn net.Conn, r *request.Request) {
		if _, err := os.Stat(router.StaticDirectory + r.Path); err == nil {
			value, err := os.ReadFile(router.StaticDirectory + r.Path)
			if err == nil {
				write := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(value)) + "\r\n\r\n" + string(value)
				router.Conn.Write([]byte(write))
			} else {
				fmt.Println("Error while loading file", err)
			}
		}
	})
	router.ServeStatic(staticDirectory)

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
