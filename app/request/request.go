package request

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const (
	METHOD_GET = iota
	METHOD_POST
)

type Request struct {
	buffer  []byte
	Method  int
	Path    string
	Params  map[string]string
	Headers Header
	Body    string
}

type Header struct {
	UserAgent string
	Host      string
}

func NewRequest(buffer []byte) Request {

	req := Request{buffer, 0, "", make(map[string]string, 0), Header{"", ""}, ""}

	req.getPath()
	req.parseHeader()
	req.parseBody()

	return req
}

func (request *Request) parseHeader() {
	segments := strings.Split(string(request.buffer), "\r\n\r\n")

	if len(segments) < 2 {
		panic("No headers nor body")
	}

	headers := strings.Split(segments[0], "\r\n")

	for _, value := range headers {
		line := strings.Split(value, " ")
		if len(line) < 1 {
			break
		}

		key := strings.TrimSuffix(line[0], ":")

		switch key {
		case "GET":
			request.Method = METHOD_GET
		case "POST":
			request.Method = METHOD_POST
		case "User-Agent":
			request.Headers.UserAgent = strings.Join(line[1:], " ")
		case "Host":
			request.Headers.Host = strings.Join(line[1:], " ")
		default:
			//noting
		}
	}

}

func (request *Request) parseBody() {
	segments := strings.Split(string(request.buffer), "\r\n\r\n")

	body := make([]byte, 0)

	for _, val := range body {
		if val != 0 {
			body = append(body, val)
		}
	}

	request.Body = string(body)
}

func (request *Request) getPath() string {

	sn := 0

	for sn <= len(request.buffer) {
		lineIndex := bytes.Index(request.buffer, []byte("\r\n"))

		if lineIndex == -1 {
			fmt.Println("Error in the header")
			os.Exit(1)
		}

		line := request.buffer[sn:lineIndex]

		for lineIndex != -1 {
			slices := bytes.Split(line, []byte(" "))

			if len(slices) != 3 {
				break
			}

			if string(slices[0]) == "GET" || string(slices[0]) == "POST" {
				//HTTPmethod = string(slices[0])
				request.Path = string(slices[1])
				return request.Path
			}
		}

		sn = lineIndex + 2
	}

	return request.Path
}
