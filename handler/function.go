package handler

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"router-practice/logger"
	"router-practice/model"
	"router-practice/router"
	"router-practice/variable"

	"github.com/goccy/go-json"
)

func Index(c *router.Context) {
	c.URL.Path = "/index.html"
	HandleHTML(c)
}

func Hello(c *router.Context) {
	if c.Method == "GET" {
		c.Text(http.StatusOK, "Hello world GET")
	} else if c.Method == "POST" {
		c.Text(http.StatusOK, "Hello world POST")
	}
}

func HelloParam(c *router.Context) {
	if len(c.Params) > 0 {
		c.Text(http.StatusOK, "Hello "+c.Params[0])
	} else {
		c.Text(http.StatusBadRequest, "Missing parameter")
	}
}

func GetParam(c *router.Context) {
	result := ""

	params := c.URL.Query()

	for k := range c.URL.Query() {
		result += k + "=" + params.Get(k) + "\n"
	}

	c.Text(http.StatusOK, result)
}

func PostForm(c *router.Context) {
	result := ""

	switch c.Method {
	case "GET":
		result = "Hello world GET"
	case "POST":
		c.ParseForm()
		for k := range c.Form {
			result += k + "=" + c.FormValue(k) + "\n"
		}
	}

	c.Text(http.StatusOK, result)
}

func PostJson(c *router.Context) {
	user := model.UserInfo{}

	err := json.NewDecoder(c.Body).Decode(&user)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	result, err := json.Marshal(user)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	c.Text(http.StatusOK, string(result))
}

func HandleHTML(c *router.Context) {
	var h []byte
	var err error

	if _, er := os.Stat("../html/" + path.Base(c.URL.Path)); er == nil {
		h, err = os.ReadFile("../html/" + path.Base(c.URL.Path)) // Real storage
	} else {
		h, err = variable.Content.ReadFile("html/" + path.Base(c.URL.Path)) // Embed storage
	}

	if err != nil {
		logger.Object.Warn().Err(err).Msg("StaticHTML")
	}

	h = bytes.ReplaceAll(h, []byte("#USERNAME"), []byte("Robert Garcia"))

	c.Html(http.StatusOK, h)
}

func HandleAsset(c *router.Context) {
	var h []byte
	var err error

	if _, er := os.Stat("../html" + c.URL.Path); er == nil {
		h, err = os.ReadFile("../html" + c.URL.Path) // Real storage
	} else {
		h, err = variable.Content.ReadFile("html" + c.URL.Path) // Embed storage
	}

	if err != nil {
		logger.Object.Warn().Err(err).Msg("StaticFiles")
	}

	c.File(http.StatusOK, h)
}
