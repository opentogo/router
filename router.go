package router

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

var (
	paramsKey  int
	parameters = map[string]string{}
)

type route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

// Router base struct
type Router struct {
	routes   []route
	notFound http.HandlerFunc
}

func (r Router) parse(path, requestURL string) bool {
	var (
		namedParams []string
		pattern     string
	)
	for _, part := range strings.Split(path, "/") {
		if len(part) == 0 {
			continue
		}
		if strings.Contains(part, ":") {
			namedParams = append(namedParams, strings.Split(part, ":")[1])
			pattern += "/(.+)"
			continue
		}
		pattern += "/" + part
	}
	if namedValues := regexp.MustCompile(pattern).FindStringSubmatch(requestURL); namedValues != nil {
		if namedValues[0] == requestURL {
			for i, value := range namedValues[1:] {
				parameters[namedParams[i]] = value
			}
			return true
		}
	}
	return false
}

// Handler registers handlers with the given path and method
func (r *Router) Handler(method string, path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		path:    path,
		method:  method,
		handler: handler,
	})
}

// NotFoundHandler registers a handler for "404 not found" requests
func (r *Router) NotFoundHandler(handler http.HandlerFunc) {
	r.notFound = handler
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.method != req.Method {
			continue
		}
		if route.path == req.URL.Path {
			route.handler(w, req)
			return
		}
		if handler := route.handler; handler != nil {
			if ok := r.parse(route.path, req.URL.Path); ok {
				handler(w, req.WithContext(context.WithValue(req.Context(), paramsKey, parameters)))
				return
			}
		}
	}
	if r.notFound != nil {
		r.notFound(w, req)
		return
	}
	http.NotFound(w, req)
}

// Params returns route param stored in http.request.
func Params(r *http.Request) map[string]string {
	if values, ok := r.Context().Value(paramsKey).(map[string]string); ok {
		return values
	}
	return nil
}
