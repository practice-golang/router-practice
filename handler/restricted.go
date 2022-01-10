package handler

import (
	"net/http"

	"router-practice/auth"
	"router-practice/model"
	"router-practice/router"
)

func RestrictedHello(c *router.Context) {
	authinfo := c.AuthInfo.(model.AuthInfo)
	c.Text(http.StatusOK, "Hello "+authinfo.Name.String)
}

// SignOut - Expire cookie
func SignOut(c *router.Context) {
	auth.ExpireCookie(c.ResponseWriter)
	authinfo := c.AuthInfo.(model.AuthInfo)
	c.Text(http.StatusOK, "Good bye "+authinfo.Name.String)
}
