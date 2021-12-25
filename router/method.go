package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"regexp"
	"router-practice/variable"
)

func NewApp() *App {
	app := &App{
		DefaultRoute: func(c *Context) {
			c.Text(http.StatusNotFound, "Not found")
		},
		MethodNotAllowed: func(c *Context) {
			c.Text(http.StatusNotFound, "Method not allowed")
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
	c := &Context{Request: r, ResponseWriter: w}

	b, _ := ioutil.ReadAll(c.Body)
	c.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	log := variable.Logger.Log()
	if json.Valid(b) {
		log = log.RawJSON("body", b)
	} else {
		log = log.Fields(map[string]interface{}{"body": b})
	}

	log.
		Timestamp().
		Str("method", c.Method).
		Str("path", c.URL.Path).
		Str("remote", c.RemoteAddr).
		Str("user-agent", c.UserAgent()).
		Fields(map[string]interface{}{"header": c.Request.Header}).
		Send()

	for _, rt := range a.Routes {
		if matches := rt.Pattern.FindStringSubmatch(c.URL.Path); len(matches) > 0 {

			if !rt.Methods[c.Request.Method] {
				// a.MethodNotAllowed(c)
				a.DefaultRoute(c)
				return
			}

			if len(matches) > 1 {
				c.Params = matches[1:]
			}

			rt.Handler(c)
			return
		}
	}

	a.DefaultRoute(c)
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
	StaticContent, err := fs.Sub(fs.FS(variable.Static), "static")
	if err != nil {
		variable.Logger.Fatal().Err(err).Msg("SetupStatic")
	}
	s := http.StripPrefix("/static/", http.FileServer(http.FS(StaticContent)))
	// s := http.StripPrefix("/static/", http.FileServer(http.Dir("../static")))
	StaticServer = func(c *Context) { s.ServeHTTP(c.ResponseWriter, c.Request) }
}
