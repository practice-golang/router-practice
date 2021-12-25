package main // import "router-practice"

import (
	"embed"
	"net/http"
	"router-practice/handler"
	"router-practice/logger"
	"router-practice/router"
	"router-practice/variable"
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
	a := router.NewApp()

	a.Handle(`^/hello$`, handler.Hello, "GET", "POST")
	a.Handle(`/hello/([\w\._-]+)$`, handler.HelloParam, "GET")

	a.Handle(`^/login$`, handler.Login, "GET", "POST")
	a.Handle(`^/user$`, handler.User, "POST")

	a.Handle(`/[^/]+.html`, handler.StaticHTML, "GET")
	a.Handle(`/html/*`, handler.StaticFiles, "GET")

	a.Handle(`/static/*`, router.StaticServer, "GET")

	variable.Logger.Log().Timestamp().Str("listen", uri+"\n").Send()

	err := http.ListenAndServe(uri, a)
	if err != nil {
		variable.Logger.Fatal().Err(err).Timestamp().Msg("Server start failed")
	}
}
