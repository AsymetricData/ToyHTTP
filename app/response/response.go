package response

import (
	"net"
	"strconv"
)

type Response struct {
	conn    net.Conn
	headers map[string]string
}

func NewResponse(c net.Conn) Response {
	defaultHeaders := make(map[string]string)

	defaultHeaders["HTTP/1.1"] = "200 OK"
	defaultHeaders["Content-Type:"] = "text/plain"

	return Response{c, defaultHeaders}
}

func (response *Response) SetHeader(header string, value string) {
	if header != "HTTP/1.1" {
		response.headers[header+":"] = value
	} else {
		response.headers[header] = value
	}
}

func (response *Response) SetStatus(status int) {
	switch status {
	case 404:
		response.SetHeader("HTTP/1.1", "404 Not Found")
		response.Write("")
	}
}

func (response *Response) Write(data string) {
	len := len(data)

	buffer := ""
	//

	for index, value := range response.headers {
		buffer = buffer + index + " " + value + " \r\n"
	}

	buffer = buffer + "Content-Length: " + strconv.Itoa(len) + "\r\n\r\n" + data

	response.conn.Write([]byte(buffer))
}
