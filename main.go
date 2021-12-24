package main // import "router-practice"

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"router-practice/router"
)

//go:embed html/*
var Content embed.FS

//go:embed static/*
var Static embed.FS

var staticServer router.Handler

func setupStatic() {
	StaticContent, err := fs.Sub(fs.FS(Static), "static")
	if err != nil {
		log.Fatal(err)
	}
	s := http.StripPrefix("/static/", http.FileServer(http.FS(StaticContent)))
	// s := http.StripPrefix("/static/", http.FileServer(http.Dir("../static")))
	staticServer = func(ctx *router.Context) { s.ServeHTTP(ctx.ResponseWriter, ctx.Request) }
}

func main() {
	uri := "127.0.0.1:4416"

	setupStatic()
	a := router.NewApp()

	a.Handle(`^/hello$`, func(ctx *router.Context) {
		ctx.Text(http.StatusOK, "Hello world")
	}, "GET", "POST")

	a.Handle(`/hello/([\w\._-]+)$`, func(ctx *router.Context) {
		ctx.Text(http.StatusOK, fmt.Sprintf("Hello %s", ctx.Params[0]))
	}, "GET")

	a.Handle(`/[^/]+.html`, func(ctx *router.Context) {
		var h []byte
		var err error
		filePATH := "../html/" + path.Base(ctx.URL.Path)
		if _, er := os.Stat(filePATH); er == nil {
			h, err = os.ReadFile(filePATH)
		} else {
			h, err = Content.ReadFile("html/" + path.Base(ctx.URL.Path))
		}

		if err != nil {
			log.Fatal(err)
		}

		ctx.Html(http.StatusOK, h)
	}, "GET")

	a.Handle(`/html/*`, func(ctx *router.Context) {
		var h []byte
		var err error
		filePATH := "../" + ctx.URL.Path[1:]
		if _, er := os.Stat(filePATH); er == nil {
			h, err = os.ReadFile(filePATH)
		} else {
			h, err = Content.ReadFile(ctx.URL.Path[1:])
		}

		if err != nil {
			log.Println(err)
		}

		ctx.Text(http.StatusOK, string(h))
	}, "GET")

	a.Handle(`/static/*`, staticServer, "GET")

	fmt.Println(uri)

	err := http.ListenAndServe(uri, a)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
