package handler

import (
	"log"
	"net/http"
	"os"
	"path"
	"router-practice/router"
	"router-practice/value"
)

func Hello(ctx *router.Context) {
	ctx.Text(http.StatusOK, "Hello world")
}

func HelloParam(ctx *router.Context) {
	ctx.Text(http.StatusOK, "Hello "+ctx.Params[0])
}

func StaticHTML(ctx *router.Context) {
	var h []byte
	var err error
	filePATH := "../html/" + path.Base(ctx.URL.Path)
	if _, er := os.Stat(filePATH); er == nil {
		h, err = os.ReadFile(filePATH)
	} else {
		h, err = value.Content.ReadFile("html/" + path.Base(ctx.URL.Path))
	}

	if err != nil {
		log.Fatal(err)
	}

	ctx.Html(http.StatusOK, h)
}

func StaticFiles(ctx *router.Context) {
	var h []byte
	var err error
	filePATH := "../html/" + path.Base(ctx.URL.Path)
	if _, er := os.Stat(filePATH); er == nil {
		h, err = os.ReadFile(filePATH)
	} else {
		h, err = value.Content.ReadFile("html/" + path.Base(ctx.URL.Path))
	}

	if err != nil {
		log.Fatal(err)
	}

	ctx.Text(http.StatusOK, string(h))
}
