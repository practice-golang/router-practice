package handler

import (
	"bytes"
	"net/http"
	"os"
	"router-practice/logging"
	"router-practice/model"
	"router-practice/router"
	"router-practice/util"
	"router-practice/wsock"

	"github.com/goccy/go-json"
)

func Index(c *router.Context) {
	c.URL.Path = "/index.html"
	HandleHTML(c)
}

func HealthCheck(c *router.Context) {
	c.Text(http.StatusOK, "Ok")
}

func Hello(c *router.Context) {
	switch c.Method {
	case http.MethodGet:
		c.Text(http.StatusOK, "Hello world GET")
	case http.MethodPost:
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
	case http.MethodGet:
		result = "Hello world GET"
	case http.MethodPost:
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

	c.Json(http.StatusOK, user)
}

func HandleHTML(c *router.Context) {
	var h []byte
	var err error

	// If the file exists in the real storage, read real instead of embed.
	storePath := StoreRoot + c.URL.Path // Real storage
	embedPath := EmbedRoot + c.URL.Path // Embed storage
	switch true {
	case util.CheckFileExists(storePath, false):
		h, err = os.ReadFile(storePath)
	case util.CheckFileExists(embedPath, true):
		h, err = router.Content.ReadFile(embedPath)
	default:
		c.Text(http.StatusNotFound, "Not found")
		return
	}

	if err != nil {
		logging.Object.Warn().Err(err).Msg("HandleHTML")
	}

	h = bytes.ReplaceAll(h, []byte("#USERNAME"), []byte("Robert Garcia"))

	c.Html(http.StatusOK, h)
}

func HandleAsset(c *router.Context) {
	var h []byte
	var err error

	// If the file exists in the real storage, read real instead of embed.
	storePath := StoreRoot + c.URL.Path // Real storage
	embedPath := EmbedRoot + c.URL.Path // Embed storage
	switch true {
	case util.CheckFileExists(storePath, false):
		h, err = os.ReadFile(storePath)
	case util.CheckFileExists(embedPath, true):
		h, err = router.Content.ReadFile(embedPath)
	default:
		c.Text(http.StatusNotFound, "Not found")
		return
	}

	if err != nil {
		logging.Object.Warn().Err(err).Msg("HandleAsset")
	}

	c.File(http.StatusOK, h)
}

func HandleWebsocketEcho(c *router.Context) {
	wsock.WebSocketEcho(c.ResponseWriter, c.Request)
}

func HandleWebsocketChat(c *router.Context) {
	wsock.WebSocketChat(c.ResponseWriter, c.Request)
}
