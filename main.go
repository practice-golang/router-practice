package main // import "router-practice"

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"router-practice/handler"
	"router-practice/router"
	"router-practice/value"
)

//go:embed html/*
var Content embed.FS

//go:embed static/*
var Static embed.FS

func main() {
	uri := "127.0.0.1:4416"

	value.Content = Content
	value.Static = Static

	router.SetupStatic()
	a := router.NewApp()

	a.Handle(`^/hello$`, handler.Hello, "GET", "POST")
	a.Handle(`/hello/([\w\._-]+)$`, handler.HelloParam, "GET")

	a.Handle(`/[^/]+.html`, handler.StaticHTML, "GET")
	a.Handle(`/html/*`, handler.StaticFiles, "GET")

	a.Handle(`/static/*`, router.StaticServer, "GET")

	fmt.Println(uri)

	err := http.ListenAndServe(uri, a)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
