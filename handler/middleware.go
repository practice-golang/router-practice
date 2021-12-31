package handler

import (
	"errors"
	"log"
	"net/http"
	"router-practice/router"
)

func HelloMiddleware(c *router.Context) error {
	log.Println("Hello middleware Sorry error")

	c.Text(http.StatusInternalServerError, "Middle ware test error")

	return errors.New("hello middleware error")
}
