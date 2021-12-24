package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/wesovilabs/koazee"
)

func NewApp() *App {
	app := &App{
		DefaultRoute: func(ctx *Context) {
			ctx.Text(http.StatusNotFound, "Not found")
		},
	}

	return app
}

func (a *App) Handle(pattern string, handler Handler, methods ...string) {
	re := regexp.MustCompile(pattern)
	route := Route{Pattern: re, Handler: handler, Methods: methods}

	a.Routes = append(a.Routes, route)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{Request: r, ResponseWriter: w}

	for _, rt := range a.Routes {
		if matches := rt.Pattern.FindStringSubmatch(ctx.URL.Path); len(matches) > 0 {
			log.Println("Method: ", rt.Methods, ctx.Request.Method)

			isMethodsExists, _ := koazee.StreamOf(rt.Methods).Contains(ctx.Request.Method)
			if !isMethodsExists {
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
