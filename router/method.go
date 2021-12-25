package router

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"regexp"
	"router-practice/value"
)

func NewApp() *App {
	app := &App{
		DefaultRoute: func(ctx *Context) {
			ctx.Text(http.StatusNotFound, "Not found")
		},
		MethodNotAllowed: func(ctx *Context) {
			ctx.Text(http.StatusNotFound, "Method not allowed")
		},
	}

	return app
}

func (a *App) Handle(pattern string, handler Handler, methods ...string) {
	re := regexp.MustCompile(pattern)
	m := Methods{}

	for _, method := range methods {
		switch method {
		case "GET":
			m["GET"] = true
		case "HEAD":
			m["HEAD"] = true
		case "POST":
			m["POST"] = true
		case "PUT":
			m["PUT"] = true
		case "PATCH":
			m["PATCH"] = true
		case "DELETE":
			m["DELETE"] = true
		case "*":
			m["GET"] = true
			m["HEAD"] = true
			m["POST"] = true
			m["PUT"] = true
			m["PATCH"] = true
			m["DELETE"] = true
		}
	}

	route := Route{Pattern: re, Handler: handler, Methods: m}

	a.Routes = append(a.Routes, route)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{Request: r, ResponseWriter: w}

	for _, rt := range a.Routes {
		if matches := rt.Pattern.FindStringSubmatch(ctx.URL.Path); len(matches) > 0 {

			if !rt.Methods[ctx.Request.Method] {
				// a.MethodNotAllowed(ctx)
				a.DefaultRoute(ctx)
				return
			}

			if len(matches) > 1 {
				ctx.Params = matches[1:]
			}

			rt.Handler(ctx)
			return
		}
	}

	a.DefaultRoute(ctx)
}

func (c *Context) Text(code int, body string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/plain")
	c.WriteHeader(code)

	io.WriteString(c.ResponseWriter, fmt.Sprintf("%s\n", body))
}

// func (c *Context) Html(code int, body string) {
func (c *Context) Html(code int, body []byte) {
	c.ResponseWriter.Header().Set("Content-Type", "text/html")
	c.WriteHeader(code)

	io.Writer(c.ResponseWriter).Write(body)
	// io.WriteString(c.ResponseWriter, fmt.Sprintf("%s\n", body))
}

var StaticServer Handler

func SetupStatic() {
	StaticContent, err := fs.Sub(fs.FS(value.Static), "static")
	if err != nil {
		log.Fatal(err)
	}
	s := http.StripPrefix("/static/", http.FileServer(http.FS(StaticContent)))
	// s := http.StripPrefix("/static/", http.FileServer(http.Dir("../static")))
	StaticServer = func(ctx *Context) { s.ServeHTTP(ctx.ResponseWriter, ctx.Request) }
}
