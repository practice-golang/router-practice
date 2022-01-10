package main

import (
	"os"
	"time"

	"router-practice/auth"
	"router-practice/handler"
	"router-practice/logging"
	"router-practice/router"
	"router-practice/wsock"

	"github.com/rs/cors"
)

func setupKey() {
	auth.Secret = "practice-golang/router-practice secret"
	err := auth.GenerateKey()

	if err != nil {
		panic(err)
	}
}

func setupLogger() {
	logging.SetupLogger()

	go func() {
		now := time.Now()
		zone, i := now.Zone()
		nextDay := now.AddDate(0, 0, 1).In(time.FixedZone(zone, i))
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		restTimeNextDay := time.Until(nextDay)
		time.Sleep(restTimeNextDay)
		for {
			if time.Now().Format("15") == "00" {
				logging.RenewLogger()
				time.Sleep(24 * time.Hour)
			} else {
				time.Sleep(time.Second)
			}
		}
	}()
}

func setupRouter() {
	// router.StaticPath = "./static"
	router.StaticPath = StaticPath
	router.Content = Content
	router.EmbedStatic = EmbedStatic

	router.SetupStaticServer()

	r := router.New()

	/* API */
	g := r.Group(`^/api`)
	g.Handle(`/?$`, handler.HealthCheck, "GET")
	g.Handle(`/hello$`, handler.Hello, "GET")
	g.Handle(`/signin$`, handler.SigninAPI, "POST")

	/* File & Directory */
	g.Handle(`/dir/list$`, handler.HandleGetDir, "POST")

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
	r.Handle(`/embed/*`, router.EmbedStaticServer, "GET")

	/* Websocket - /ws.html */
	r.Handle(`/ws-echo`, handler.HandleWebsocketEcho, "GET")
	r.Handle(`/ws-chat`, handler.HandleWebsocketChat, "GET")

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

func doSetup() {
	_ = os.Mkdir(StaticPath, os.ModePerm)

	setupKey()
	setupLogger()
	setupRouter()

	wsock.InitWebSocketChat()
}
