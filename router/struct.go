package router

import (
	"net/http"
	"regexp"
)

type Handler func(*Context)

type Methods map[string]bool

type Route struct {
	Pattern        *regexp.Regexp
	Handler        Handler
	Methods Methods
}

type App struct {
	Routes           []Route
	DefaultRoute     Handler
	MethodNotAllowed Handler
}

type Context struct {
	http.ResponseWriter
	*http.Request
	Params []string
}
