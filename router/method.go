package router

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strings"

	"router-practice/logging"
)

var StaticServer Handler
var EmbedStaticServer Handler

func New(middleware ...Middleware) *App {
	app := &App{
		DefaultRoute: func(c *Context) {
			c.Text(http.StatusNotFound, "Not found")
		},
		MethodNotAllowed: func(c *Context) {
			c.Text(http.StatusNotFound, "Method not allowed")
		},
		Middlewares: middleware,
	}

	return app
}

func (a *App) Group(prefix string, middleware ...Middleware) *RouteGroup {
	group := &RouteGroup{
		App:         a,
		Prefix:      prefix,
		Middlewares: middleware,
	}

	return group
}

func (g *RouteGroup) Handle(pattern string, handler Handler, methods ...string) {
	appMiddlewares := g.App.Middlewares
	g.App.Middlewares = append(g.App.Middlewares, g.Middlewares...)
	g.App.Handle(g.Prefix+pattern, handler, methods...)
	g.App.Middlewares = appMiddlewares
}

func (a *App) Handle(pattern string, handler Handler, methods ...string) {
	re := regexp.MustCompile(pattern)
	m := Methods{}

	for _, method := range methods {
		switch method {
		case "*":
			for _, method := range AllMethods {
				m[method] = true
			}
		default:
			m[strings.ToUpper(method)] = true
		}
	}

	route := Route{Pattern: re, Handler: handler, Methods: m, Middlewares: a.Middlewares}

	a.Routes = append(a.Routes, route)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{Request: r, ResponseWriter: w}

	b, _ := ioutil.ReadAll(c.Body)
	c.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	logger := logging.Object.Log()
	if json.Valid(b) {
		bc := new(bytes.Buffer)
		json.Compact(bc, b)
		logger = logger.RawJSON("body", bc.Bytes())
	} else {
		logger = logger.Fields(map[string]interface{}{"body": bytes.ReplaceAll(b, []byte("\n"), []byte(""))})
	}

	logger.Timestamp().
		Str("method", c.Method).
		Str("path", c.URL.Path).
		Str("remote", c.RemoteAddr).
		Str("user-agent", c.UserAgent()).
		Fields(map[string]interface{}{"header": c.Request.Header}).
		Send()

	for _, rt := range a.Routes {
		matches := rt.Pattern.FindStringSubmatch(c.URL.Path)
		if len(matches) > 0 {
			// log.Println("Route path regex:", rt.Pattern.String(), c.URL.Path, matches, len(matches))
			if !rt.Methods[c.Method] {
				// a.MethodNotAllowed(c)
				a.DefaultRoute(c)
				return
			}

			if len(matches) > 1 {
				c.Params = matches[1:]
			}

			for _, m := range rt.Middlewares {
				rt.Handler = m(rt.Handler)
			}

			rt.Handler(c)
			return
		}
	}

	a.DefaultRoute(c)
}

func (c *Context) Text(code int, body string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	c.WriteHeader(code)

	c.ResponseWriter.Write([]byte(body))
}

// func (c *Context) Json(code int, body string) {
func (c *Context) Json(code int, body interface{}) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	c.WriteHeader(code)

	result, err := json.Marshal(body)
	if err != nil {
		log.Println("Json error:", err)
		return
	}

	c.ResponseWriter.Write(result)
}

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

// SetupStaticServer - Serving internal `embedded` static files
func SetupStaticServer() {
	var err error

	EmbedContent, err = fs.Sub(fs.FS(EmbedStatic), EmbedPath)
	if err != nil {
		logging.Object.Warn().Err(err).Msg("SetupStatic")
	}

	e := http.StripPrefix("/embed/", http.FileServer(http.FS(EmbedContent))) // embed storage
	EmbedStaticServer = func(c *Context) { e.ServeHTTP(c.ResponseWriter, c.Request) }

	s := http.StripPrefix("/static/", http.FileServer(http.Dir(StaticPath))) // real storage
	StaticServer = func(c *Context) { s.ServeHTTP(c.ResponseWriter, c.Request) }
}
