package main // import "router-practice"

import (
	"embed"
	"net/http"
	"router-practice/handler"
	"router-practice/logger"
	"router-practice/router"
	"router-practice/variable"

	"github.com/rs/cors"
)

//go:embed html/*
var Content embed.FS

//go:embed static/*
var Static embed.FS

func main() {
	uri := "127.0.0.1:4416"

	variable.Content = Content
	variable.Static = Static

	logger.SetupLogger()

	router.SetupStatic()
	a := router.New()

	all := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	a.Handle(`^[/|]$`, handler.Index, "GET")

	a.Handle(`^/hello$`, handler.Hello, "GET", "POST")
	a.Handle(`/hello/([\w\._-]+)$`, handler.HelloParam, "GET")

	a.Handle(`/get-param$`, handler.GetParam, "GET")

	a.Handle(`^/post-form$`, handler.PostForm, "GET", "POST")
	a.Handle(`^/post-json$`, handler.PostJson, all...)

	a.Handle(`/[^/]+.html`, handler.HandleHTML, "GET")
	a.Handle(`^/.*.[css|js|map]$`, handler.HandleAsset, "GET")

	a.Handle(`/static/*`, router.StaticServer, "GET")

	handler := cors.Default().Handler(a)
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:4416"},
	// 	AllowedMethods:   []string{"GET"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	Debug:            false,
	// })
	// handler := c.Handler(router)

	logger.Object.Log().Timestamp().Str("listen", uri+"\n").Send()
	println("Listen", uri)

	err := http.ListenAndServe(uri, handler)
	if err != nil {
		logger.Object.Warn().Err(err).Timestamp().Msg("Server start failed")
	}
}
