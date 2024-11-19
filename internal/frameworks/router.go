package frameworks

import (
	"net/http"
	"strings"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Handle(method, pattern string, handler http.HandlerFunc) {
	route := Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	}
	r.routes = append(r.routes, route)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method == route.Method && strings.HasPrefix(req.URL.Path, route.Pattern) {
			route.Handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}
