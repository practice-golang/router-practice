package router

import (
	"log"
	"net/http"
)

func HelloMiddleware(next Handler) Handler {
	return func(c *Context) {
		log.Println("Hello middleware")

		c.Params = append(c.Params, "earth")
		c.Text(http.StatusOK, "Earth")

		next(c)
	}
}

func ByeMiddleware(next Handler) Handler {
	return func(c *Context) {
		log.Println("Bye middleware")

		c.Params = append(c.Params, "universe")
		c.Text(http.StatusOK, "Universe")

		next(c)
	}
}
