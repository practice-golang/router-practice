package main // import "router-practice"

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"router-practice/router"
)

type CustomHandler struct {
	http.Handler
}

//go:embed html/*
var Content embed.FS

func HelloWorld(c *router.Context) {
	// log.Println("WTF???")
	c.Text(http.StatusOK, fmt.Sprintf("Hello %s", c.Params[0]))
}

func HelloMux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	uri := "127.0.0.1:4416"

	handler := CustomHandler{
		Handler: http.HandlerFunc(HelloMux),
	}

	b := http.NewServeMux()

	b.Handle("/", handler)

	a := router.NewApp()

	a.Handle(`^/hello$`, func(ctx *router.Context) {
		ctx.Text(http.StatusOK, "Hello world")
	})

	a.Handle(`/hello/([\w\._-]+)$`, HelloWorld)

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
