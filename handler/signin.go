package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"router-practice/auth"
	"router-practice/model"
	"router-practice/router"

	"gopkg.in/guregu/null.v4"
)

func Signin(c *router.Context) {
	var err error

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Json(http.StatusBadRequest, err.Error())
	}

	signin := model.SignIn{}
	err = json.Unmarshal(b, &signin)
	if err != nil {
		c.Json(http.StatusInternalServerError, err.Error())
	}

	authinfo := model.AuthInfo{
		Name:     null.NewString(signin.Name.String, true),
		IpAddr:   null.NewString(c.RemoteAddr, true),
		Platform: null.NewString("", true),
		Duration: null.NewInt(60*60*24*7, true),
		// Duration: null.NewInt(10, true), // 10 seconds test
	}

	// auth.SetupCookieToken(c.ResponseWriter, authinfo)
	auth.SetCookieSession(c, authinfo)

	c.Json(http.StatusOK, "Signin success")
}

func Login(c *router.Context) {
	failBody := `<meta http-equiv="refresh" content="2; url=/"></meta>`

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		c.Html(http.StatusBadRequest, []byte(failBody+`Missing parameter`))
		return
	}

	authinfo := model.AuthInfo{
		Name:     null.NewString(username, true),
		IpAddr:   null.NewString(c.RemoteAddr, true),
		Platform: null.NewString("", true),
		Duration: null.NewInt(60*60*24*7, true),
		// Duration: null.NewInt(10, true), // 10 seconds test
	}

	auth.SetCookieSession(c, authinfo)

	// c.Json(http.StatusOK, "Signin success")

	destination := "/"
	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=`+destination+`"></meta>`))
}

func SigninAPI(c *router.Context) {
	var err error

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Json(http.StatusBadRequest, err.Error())
	}

	signin := model.SignIn{}
	err = json.Unmarshal(b, &signin)
	if err != nil {
		c.Json(http.StatusInternalServerError, err.Error())
	}

	authinfo := model.AuthInfo{
		Name:     null.NewString(signin.Name.String, true),
		IpAddr:   null.NewString(c.RemoteAddr, true),
		Platform: null.NewString("", true),
		Duration: null.NewInt(60*60*24*7, true),
		// Duration: null.NewInt(10, true), // 10 seconds test
	}

	token, err := auth.GenerateToken(authinfo)
	if err != nil {
		c.Json(http.StatusInternalServerError, err.Error())
	}

	result := map[string]string{
		"token": token,
		"msg":   "Signin success",
	}

	c.Json(http.StatusOK, result)
}
