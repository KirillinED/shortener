package router

import (
	"fmt"
	"net/http"
	"regexp"
)

func NewRouter() *Router {
	return new(Router)
}

type Router struct {
	Routes []Route
}

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.Routes {
		if r.ResolveRoute(route, req) {
			fmt.Println(route.Method, route.Path)
			route.Handler(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

func (r *Router) ResolveRoute(route Route, req *http.Request) bool {
	if route.Method != req.Method {
		return false
	}

	if route.Path == req.URL.Path {
		return true
	}

	matched, err := regexp.Match(route.Path, []byte(req.URL.Path))
	if err != nil {
		return false
	}

	return matched
}

func (r *Router) Handle(method string, path string, handler http.HandlerFunc) {
	route := newRoute(method, path, handler)

	r.Routes = append(r.Routes, route)
}

func newRoute(method string, path string, handler http.HandlerFunc) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
