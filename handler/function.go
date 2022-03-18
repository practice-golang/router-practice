package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"router-practice/fd"
	"router-practice/logging"
	"router-practice/model"
	"router-practice/router"
	"router-practice/util"
	"router-practice/wsock"

	"gopkg.in/guregu/null.v4"
	// "github.com/goccy/go-json"
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

		// Parameter transfer test
		// c.Text(http.StatusOK, "Hi, "+c.Params[0]+" / Context value: "+c.Context().Value(router.ContextKey("say")).(string))
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

func HandleGetDir(c *router.Context) {
	path := model.FilePath{}
	result := model.FileList{}

	err := json.NewDecoder(c.Body).Decode(&path)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	f, err := os.Stat(path.Path.String)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	dir := path.Path.String
	if f.IsDir() {
		dir = path.Path.String + "/"
	}

	dir = filepath.Dir(dir)
	absPath, err := filepath.Abs(dir)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	sort := 0
	switch path.Sort.String {
	case "name":
		sort = fd.NAME
	case "size":
		sort = fd.SIZE
	case "time":
		sort = fd.TIME
	default:
		sort = fd.NAME
	}

	order := 0
	switch path.Order.String {
	case "asc":
		order = fd.ASC
	case "desc":
		order = fd.DESC
	default:
		order = fd.ASC
	}

	result.Path = null.StringFrom(dir)
	result.FullPath = null.StringFrom(absPath)

	files, err := fd.Dir(absPath, sort, order)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	for _, file := range files {
		result.Files = append(result.Files, model.FileInfo{
			Name:     null.StringFrom(file.Name()),
			Size:     null.IntFrom(file.Size()),
			DateTime: null.StringFrom(file.ModTime().Format("2006-01-02 15:04:05")),
			IsDir:    null.BoolFrom(file.IsDir()),
		})
	}

	c.Json(http.StatusOK, result)
}
