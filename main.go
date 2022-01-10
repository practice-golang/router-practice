package main // import "router-practice"

import (
	"embed"
	"net/http"
	"os"

	"router-practice/logging"
)

//go:embed html/*
var Content embed.FS

//go:embed embed/*
var EmbedStatic embed.FS

var StaticPath = "../static"

var (
	Address       string = "localhost"
	Port          string = "4416"
	ServerHandler http.Handler
)

func main() {
	envPORT := os.Getenv("PORT")

	if envPORT == "" {
		envPORT = "4416"
	}

	Port = envPORT

	uri := Address + ":" + Port

	doSetup()

	logging.Object.Log().Timestamp().Str("listen", Address+"\n").Send()
	println("Listen", uri)

	err := http.ListenAndServe(uri, ServerHandler)
	if err != nil {
		logging.Object.Warn().Err(err).Timestamp().Msg("Server start failed")
	}
}
