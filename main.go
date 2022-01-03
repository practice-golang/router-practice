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
	Port          string = "4416"
	Address       string = "localhost"
	ServerHandler http.Handler
)

func main() {
	uri := Address + ":" + Port

	setupLogger()
	setupKey()
	setupRouter()

	logging.Object.Log().Timestamp().Str("listen", Address+"\n").Send()
	println("Listen", uri)

	err := http.ListenAndServe(uri, ServerHandler)
	if err != nil {
		logging.Object.Warn().Err(err).Timestamp().Msg("Server start failed")
	}
}
