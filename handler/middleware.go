package handler

import (
	"errors"
	"log"
	"net/http"
	"router-practice/auth"
	"router-practice/router"
)

func HelloMiddleware(c *router.Context) error {
	c.Text(http.StatusInternalServerError, "Middle ware test error")

	return errors.New("hello middleware error")
}

func AuthMiddleware(c *router.Context) error {
	claim, err := auth.GetClaim(*c.Request)
	if err != nil {
		auth.ExpireCookie(c.ResponseWriter)

		log.Println("AuthMiddleware:", err)
		c.Text(http.StatusUnauthorized, "Auth error")

		return err
	}

	c.AuthInfo = claim
	// c.Params = append(c.Params, claim.Name.String)

	return nil
}
