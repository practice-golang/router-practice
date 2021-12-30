package main // import "router-practice"

import (
	"embed"
	"net/http"
	"router-practice/handler"
	"router-practice/logging"
	"router-practice/router"

	"github.com/rs/cors"
)

//go:embed html/*
var Content embed.FS

//go:embed static/*
var Static embed.FS

func main() {
	uri := "localhost:4416"

	logging.SetupLogger()

	router.Content = Content
	router.Static = Static

	router.SetupStaticServer()

	r := router.New()

	allMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

	g := r.Group(`^/api`)
	g.Handle(`/?$`, handler.HealthCheck, "GET")
	g.Handle(`/hello$`, handler.Hello, "GET")

	r.Handle(`^/?$`, handler.Index, "GET")

	r.Handle(`^/hello$`, handler.Hello, "GET", "POST")
	r.Handle(`/hello/[\p{L}\d_]+$`, handler.HelloParam, "GET")

	r.Handle(`/get-param$`, handler.GetParam, "GET")
	r.Handle(`^/post-form$`, handler.PostForm, "GET", "POST")
	r.Handle(`^/post-json$`, handler.PostJson, allMethods...)

	r.Handle(`^/[^/]+.html$`, handler.HandleHTML, "GET")
	r.Handle(`^/assets/css/[^/]+.html$`, handler.HandleHTML, "GET")
	r.Handle(`^/assets/.*[css|js|map|woff|woff2]$`, handler.HandleAsset, "GET")

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
