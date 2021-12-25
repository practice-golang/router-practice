package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"regexp"
	"router-practice/logger"
	"router-practice/variable"
)

func New() *App {
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
		case "*":
			m["GET"] = true
			m["HEAD"] = true
			m["POST"] = true
			m["PUT"] = true
			m["PATCH"] = true
			m["DELETE"] = true
		default:
			m[method] = true
		}
	}

	route := Route{Pattern: re, Handler: handler, Methods: m}

	a.Routes = append(a.Routes, route)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{Request: r, ResponseWriter: w}

	b, _ := ioutil.ReadAll(c.Body)
	c.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	logging := logger.Object.Log()
	if json.Valid(b) {
		logging = logging.RawJSON("body", b)
	} else {
		logging = logging.Fields(map[string]interface{}{"body": b})
	}

	logging.Timestamp().
		Str("method", c.Method).
		Str("path", c.URL.Path).
		Str("remote", c.RemoteAddr).
		Str("user-agent", c.UserAgent()).
		Fields(map[string]interface{}{"header": c.Request.Header}).
		Send()

	for _, rt := range a.Routes {
		if matches := rt.Pattern.FindStringSubmatch(c.URL.Path); len(matches) > 0 {
			log.Println(rt.Pattern.String(), c.URL.Path)

			if !rt.Methods[c.Method] {
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

	c.ResponseWriter.Write(body)
}

func (c *Context) File(code int, body []byte) {
	c.ResponseWriter.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(c.URL.Path)))
	c.WriteHeader(code)

	c.ResponseWriter.Write(body)
}

var StaticServer Handler

func SetupStatic() {
	StaticContent, err := fs.Sub(fs.FS(variable.Static), "static")
	if err != nil {
		logger.Object.Warn().Err(err).Msg("SetupStatic")
	}
	s := http.StripPrefix("/static/", http.FileServer(http.FS(StaticContent)))
	// s := http.StripPrefix("/static/", http.FileServer(http.Dir("../static")))
	StaticServer = func(c *Context) { s.ServeHTTP(c.ResponseWriter, c.Request) }
}
