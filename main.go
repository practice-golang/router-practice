package main // import "router-practice"

import (
	"embed"
	"net/http"
	"router-practice/logging"
)

//go:embed html/*
var Content embed.FS

//go:embed static/*
var Static embed.FS

var (
	Uri           string = "localhost:4416"
	ServerHandler http.Handler
)

func main() {

	setupRouter()

	logging.Object.Log().Timestamp().Str("listen", Uri+"\n").Send()
	println("Listen", Uri)

	err := http.ListenAndServe(Uri, ServerHandler)
	if err != nil {
		logging.Object.Warn().Err(err).Timestamp().Msg("Server start failed")
	}
}
