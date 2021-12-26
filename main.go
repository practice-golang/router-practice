package main // import "router-practice"

import (
	"embed"
	"net/http"
	"router-practice/handler"
	"router-practice/logging"
	"router-practice/router"
	"router-practice/variable"

	"github.com/rs/cors"
)

//go:embed html/*
var Content embed.FS

//go:embed static/*
var Static embed.FS

func main() {
	uri := "localhost:4416"

	variable.Content = Content
	variable.Static = Static

	logging.SetupLogger()

	router.SetupStaticServer()
	r := router.New()

	allMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	r.Handle(`^[/|]$`, handler.Index, "GET")

	r.Handle(`^/hello$`, handler.Hello, "GET", "POST")
	r.Handle(`/hello/([\w\._-]+)$`, handler.HelloParam, "GET")

	r.Handle(`/get-param$`, handler.GetParam, "GET")

	r.Handle(`^/post-form$`, handler.PostForm, "GET", "POST")
	r.Handle(`^/post-json$`, handler.PostJson, allMethods...)

	r.Handle(`/[^/]+.html`, handler.HandleHTML, "GET")
	r.Handle(`^/.*.[css|js|map]$`, handler.HandleAsset, "GET")

	r.Handle(`/static/*`, router.StaticServer, "GET")

	serverHandler := cors.Default().Handler(r)
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://"+listen},
	// 	AllowedMethods:   []string{"GET"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	Debug:            false,
	// })
	// serverHandler := c.Handler(r)

	logging.Object.Log().Timestamp().Str("listen", uri+"\n").Send()
	println("Listen", uri)

	err := http.ListenAndServe(uri, serverHandler)
	if err != nil {
		logging.Object.Warn().Err(err).Timestamp().Msg("Server start failed")
	}
}
