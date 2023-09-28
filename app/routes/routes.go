package routes

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
)

type Route struct {
	Path    string
	Params  []string
	Handler func(conn net.Conn, request *request.Request)
}

type Router struct {
	Routes          []Route
	BasePath        string
	Conn            net.Conn
	StaticDirectory string
}

func NewRouter(basePath string, conn net.Conn) Router {
	return Router{make([]Route, 0), basePath, conn, ""}
}

func (router *Router) Handle(path string, handler func(conn net.Conn, r *request.Request)) {
	segments := strings.Split(path, "/")
	params := make([]string, 0)
	for index := range segments {
		cur := segments[index]

		if strings.HasPrefix(cur, "{") && strings.HasSuffix(cur, "}") {
			params = append(params, cur[1:len(cur)-1])
			path = strings.Replace(path, "/"+cur, "", 1)
		}
	}

	router.Routes = append(router.Routes, Route{path, params, handler})
}

func (router *Router) ServeStatic(directory string) {
	wd, _ := os.Getwd()
	router.StaticDirectory = wd + "/" + directory + "/"
}

func (router *Router) Get(r *request.Request) error {

	for _, route := range router.Routes {
		s := route.match(r)
		if s {
			segments := strings.Split(r.Path, "/")
			segments = segments[2:]

			if len(route.Params) >= 1 {
				for index, value := range segments {
					r.Params[route.Params[index]] = value
				}
			}
			route.Handler(router.Conn, r)
			return nil
		}
	}

	return errors.New("no route found")
}

func (route *Route) match(request *request.Request) bool {

	rs := strings.Split(request.Path, "/")

	if len(rs)-2 == len(route.Params) {
		return "/"+rs[1] == route.Path
	}

	return false
}
