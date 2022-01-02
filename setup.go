package main

import (
	"router-practice/auth"
	"router-practice/handler"
	"router-practice/logging"
	"router-practice/router"

	"github.com/rs/cors"
)

func setupKey() {
	auth.Secret = "router-practice secret key"
	auth.GenerateKey()
}

func setupRouter() {
	logging.SetupLogger()

	router.Content = Content
	router.Static = Static

	router.SetupStaticServer()

	r := router.New()

	/* API */
	g := r.Group(`^/api`)
	g.Handle(`/?$`, handler.HealthCheck, "GET")
	g.Handle(`/hello$`, handler.Hello, "GET")
	g.Handle(`/signin$`, handler.SigninAPI, "POST")

	/* Group */
	gh := r.Group(`^/hello`)
	gh.Handle(`$`, handler.Hello, "GET", "POST")
	gh.Handle(`/([\p{L}\d_]+)$`, handler.HelloParam, "GET")

	/* Middleware */
	gm := r.Group(``, handler.HelloMiddleware)
	gm.Handle(`/hi/([\p{L}\d_]+)$`, handler.HelloParam, "GET")

	/* Restricted - Cookie */
	r.Handle(`^/signin$`, handler.Signin, "POST")
	gr := r.Group(``, handler.AuthMiddleware)
	gr.Handle(`^/restricted$`, handler.RestrictedHello, "GET")
	gr.Handle(`^/signout$`, handler.SignOut, "GET")

	/* Restricted - Header */
	ga := r.Group(`^/api`, handler.AuthApiMiddleware)
	ga.Handle(`/restricted$`, handler.RestrictedHello, "GET")

	/* HTML */
	r.Handle(`^/?$`, handler.Index, "GET")

	r.Handle(`/get-param$`, handler.GetParam, "GET")
	r.Handle(`^/post-form$`, handler.PostForm, "GET", "POST")
	r.Handle(`^/post-json$`, handler.PostJson, router.AllMethods...)

	r.Handle(`^/[^/]+.html$`, handler.HandleHTML, "GET")
	r.Handle(`^/assets/css/[^/]+.html$`, handler.HandleHTML, "GET")
	r.Handle(`^/assets/.*[css|js|map|woff|woff2]$`, handler.HandleAsset, "GET")

	/* Static */
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
