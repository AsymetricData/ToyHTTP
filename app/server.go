package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
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

	if err != nil {
		//fmt.Println("Error while reading Conn ", err)
		return
	}
	//fmt.Println("Handled new data : ", n)

	router := routes.NewRouter("/", conn)
	router.Get("/", func(resp *response.Response, r *request.Request) {
		resp.Write("")
	})
	router.Get("/echo/{val}/{value}", func(resp *response.Response, r *request.Request) {
		value := r.Params["val"] + "/" + r.Params["value"]
		resp.Write(value)
	})
	router.Get("/echo/{value}", func(resp *response.Response, r *request.Request) {
		value := r.Params["value"]
		resp.Write(value)
	})
	router.Get("/user-agent", func(resp *response.Response, r *request.Request) {
		value := r.Headers.UserAgent
		resp.Write(value)
	})
	router.Get("/files/{value}", func(resp *response.Response, r *request.Request) {
		path := router.StaticDirectory + r.Params["value"]
		fmt.Println(path)
		if _, err := os.Stat(path); err == nil {
			value, err := os.ReadFile(path)
			if err == nil {
				resp.SetHeader("Content-Type", "application/octet-stream")
				resp.Write(string(value))
			} else {
				fmt.Println("Error while loading file", err)
			}
		} else {
			/* write := "HTTP/1.1 404 Not Found \r\n" */
			resp.SetStatus(404)
		}
	})
	router.Post("/files/{value}", func(resp *response.Response, r *request.Request) {
		path := router.StaticDirectory + strings.TrimPrefix(r.Params["value"], "/")
		file, err := os.Create(path)

		if file == nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		fmt.Println("Creating a file ", file.Name())

		file.WriteString(r.Body)
		resp.SetHeader("HTTP/1.1", "201 Created")
		resp.Write("")
	})
	router.ServeStatic(staticDirectory)

	r := request.NewRequest(buffer)

	err = router.Handle(&r)

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
