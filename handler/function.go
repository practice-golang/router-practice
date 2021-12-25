package handler

import (
	"log"
	"net/http"
	"os"
	"path"
	"router-practice/model"
	"router-practice/router"

	"github.com/goccy/go-json"
)

func Hello(c *router.Context) {
	if c.Method == "GET" {
		c.Text(http.StatusOK, "Hello world GET")
	} else if c.Method == "POST" {
		c.Text(http.StatusOK, "Hello world POST")
	}
}

func HelloParam(c *router.Context) {
	c.Text(http.StatusOK, "Hello "+c.Params[0])
}

func Login(c *router.Context) {
	result := ""

	if c.Method == "GET" {
		result = "Hello world GET"
	} else if c.Method == "POST" {
		result = c.FormValue("name") + "/" + c.FormValue("password") + "\n"

		c.ParseForm()
		for k, v := range c.Form {
			result += k + ":" + v[0] + "\n"
		}
	}

	c.Text(http.StatusOK, result)
}

func User(c *router.Context) {
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

func StaticHTML(c *router.Context) {
	var h []byte
	var err error
	filePATH := "../html/" + path.Base(c.URL.Path)
	if _, er := os.Stat(filePATH); er == nil {
		h, err = os.ReadFile(filePATH)
	} else {
		h, err = model.Content.ReadFile("html/" + path.Base(c.URL.Path))
	}

	if err != nil {
		log.Fatal(err)
	}

	c.Html(http.StatusOK, h)
}

func StaticFiles(c *router.Context) {
	var h []byte
	var err error
	filePATH := "../html/" + path.Base(c.URL.Path)
	if _, er := os.Stat(filePATH); er == nil {
		h, err = os.ReadFile(filePATH)
	} else {
		h, err = model.Content.ReadFile("html/" + path.Base(c.URL.Path))
	}

	if err != nil {
		log.Fatal(err)
	}

	c.Text(http.StatusOK, string(h))
}
