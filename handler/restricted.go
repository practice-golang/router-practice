package handler

import (
	"net/http"
	"router-practice/model"
	"router-practice/router"
)

func RestrictedHello(c *router.Context) {
	authinfo := c.AuthInfo.(model.AuthInfo)
	c.Text(http.StatusOK, "Hello "+authinfo.Name.String)
}
