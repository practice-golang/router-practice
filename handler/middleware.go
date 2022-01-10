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
	claim, err := auth.GetClaim(*c.Request, "cookie")
	if err != nil {
		auth.ExpireCookie(c.ResponseWriter)

		// log.Println("AuthMiddleware:", err)
		c.Text(http.StatusUnauthorized, "Auth error")

		return err
	}

	c.AuthInfo = claim
	// c.Params = append(c.Params, claim.Name.String)

	return nil
}

func AuthApiMiddleware(c *router.Context) error {
	c.Request.Header.Get("Authorization")
	claim, err := auth.GetClaim(*c.Request, "header")
	if err != nil {
		log.Println("AuthApiMiddleware:", err)
		c.Text(http.StatusUnauthorized, "Auth error")

		return err
	}

	c.AuthInfo = claim

	return nil
}
