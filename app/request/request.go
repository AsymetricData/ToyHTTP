package request

import (
	"bytes"
	"fmt"
	"os"
)

type Request struct {
	buffer []byte
	Method string
	Path   string
	Params map[string]string
}

func NewRequest(buffer []byte) Request {

	req := Request{buffer, "", "", make(map[string]string, 0)}

	req.getPath()

	return req
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
