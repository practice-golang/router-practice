package router

import (
	"net/http"
	"regexp"
)

type Handler func(*Context)

type Route struct {
	Pattern *regexp.Regexp
	Handler Handler
}

type App struct {
	Routes       []Route
	DefaultRoute Handler
}

type Context struct {
	http.ResponseWriter
	*http.Request
	Params []string
}
