package main // import "router-practice"

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"router-practice/router"
)

//go:embed html/*
var Content embed.FS

func main() {
	uri := "127.0.0.1:4416"

	a := router.NewApp()

	a.Handle(`^/hello$`, func(ctx *router.Context) {
		ctx.Text(http.StatusOK, "Hello world")
	})

	a.Handle(`/hello/([\w\._-]+)$`, func(ctx *router.Context) {
		ctx.Text(http.StatusOK, fmt.Sprintf("Hello %s", ctx.Params[0]))
	})

	a.Handle(`/*.html`, func(ctx *router.Context) {
		var h []byte
		var err error
		filePATH := "../html" + ctx.URL.Path
		if _, er := os.Stat(filePATH); er == nil {
			h, err = os.ReadFile(filePATH)
		} else {
			h, err = Content.ReadFile("html" + ctx.URL.Path)
		}

		if err != nil {
			log.Fatal(err)
		}

		ctx.Html(http.StatusOK, h)
	})

	err := http.ListenAndServe(uri, a)

	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}
