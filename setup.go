package main

import (
	"router-practice/handler"
	"router-practice/logging"
	"router-practice/router"

	"github.com/rs/cors"
)

func setupRouter() {
	logging.SetupLogger()

	router.Content = Content
	router.Static = Static

	router.SetupStaticServer()

	r := router.New()

	allMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

	g := r.Group(`^/api`)
	g.Handle(`/?$`, handler.HealthCheck, "GET")
	g.Handle(`/hello$`, handler.Hello, "GET")

	gh := r.Group(`^/hello`)
	gh.Handle(`$`, handler.Hello, "GET", "POST")
	gh.Handle(`/([\p{L}\d_]+)$`, handler.HelloParam, "GET")

	gm := r.Group(``, handler.HelloMiddleware)
	gm.Handle(`/hi/([\p{L}\d_]+)$`, handler.HelloParam, "GET")

	// HTML
	r.Handle(`^/?$`, handler.Index, "GET")

	r.Handle(`/get-param$`, handler.GetParam, "GET")
	r.Handle(`^/post-form$`, handler.PostForm, "GET", "POST")
	r.Handle(`^/post-json$`, handler.PostJson, allMethods...)

	r.Handle(`^/[^/]+.html$`, handler.HandleHTML, "GET")
	r.Handle(`^/assets/css/[^/]+.html$`, handler.HandleHTML, "GET")
	r.Handle(`^/assets/.*[css|js|map|woff|woff2]$`, handler.HandleAsset, "GET")

	r.Handle(`/static/*`, router.StaticServer, "GET")

	ServerHandler = cors.Default().Handler(r)
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://"+listen},
	// 	AllowedMethods:   []string{"GET"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	Debug:            false,
	// })
	// serverHandler := c.Handler(r)

}
