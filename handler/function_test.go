package handler

import (
	"bytes"
	"embed"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"router-practice/router"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

//go:embed embed_test/*
var fncEMBED embed.FS

func Test_Index(t *testing.T) {
	type args struct{ c *router.Context }
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_Index",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/index.html", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
		{
			name: "Test_Index_embed",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
		{
			name: "Test_Index_notfound",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeRootBackup := StoreRoot
			embedRootBackup := EmbedRoot
			compareFilePath := "../html/index.html"
			want, err := ioutil.ReadFile(compareFilePath)
			if err != nil {
				t.Error("Reference file htm not found")
			}

			if tt.name == "Test_Index_embed" {
				StoreRoot = "./not-found"
				EmbedRoot = "embed_test"
				router.Content = fncEMBED
				// tt.args.c.URL.Path = "/index.html"
				want = []byte("Hello embedded world")
			}
			if tt.name == "Test_Index_notfound" {
				StoreRoot = "./not-found"
				EmbedRoot = "not-found"
				want = []byte("Not found")
			}
			// HandleHTML(tt.args.c)
			Index(tt.args.c)

			StoreRoot = storeRootBackup
			EmbedRoot = embedRootBackup

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			want = bytes.ReplaceAll(want, []byte("#USERNAME"), []byte("Robert Garcia"))

			require.Equal(t, want, data, tt.name+" not equal"+"embed_test / "+tt.args.c.URL.Path)
		})
	}
}

func Test_HealthCheck(t *testing.T) {
	type args struct{ c *router.Context }
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_HealthCheck",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HealthCheck(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, []byte("Ok"), data, "HealthCheck not equal Ok")
		})
	}
}

func Test_Hello(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_Hello_API",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world GET"),
			},
		},
		{
			name: "Test_Hello_GET",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world GET"),
			},
		},
		{
			name: "Test_Hello_POST",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("POST", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world POST"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Hello(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "Hello not equal %v", string(data))
		})
	}
}

func Test_HelloParam(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_HelloParam",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
					Params:         []string{"test_name"},
				},
				want: []byte("Hello test_name"),
			},
		},
		{
			name: "Test_HelloParam_no_params",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Missing parameter"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HelloParam(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "HelloParam not equal")
		})
	}
}

func Test_GetParam(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_GetParam",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/get-param?hello=world", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("hello=world\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetParam(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "GetParam not equal")
		})
	}
}

func Test_PostForm(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_PostForm",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodGet, "/post-form", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world GET"),
			},
		},
		{
			name: "Test_PostForm",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodPost, "/post-form", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("hello=world\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.c.Request.Method == http.MethodPost {
				dat := url.Values{
					"hello": []string{"world"},
				}
				tt.args.c.Request.PostForm = dat
			}

			PostForm(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "PostForm not equal")
		})
	}
}

func Test_PostJson(t *testing.T) {
	jsonBody := map[string]interface{}{
		"name": "Thomas",
		"age":  "42",
	}
	body, _ := json.Marshal(jsonBody)

	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_PostJson",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodPost, "/post-json", bytes.NewReader(body)),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte(`{"name":"Thomas","age":42}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostJson(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "PostJson not equal")
		})
	}
}

func Test_HandleAsset(t *testing.T) {
	css, _ := ioutil.ReadFile("../html/assets/css/bootstrap.min.css")
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_HandleAsset",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/assets/css/bootstrap.min.css", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: css,
			},
		},
		{
			name: "Test_HandleAsset_embed",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/index.html", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello embedded world"),
			},
		},
		{
			name: "Test_HandleAsset_notfound",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/not-found", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Not found"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeRootBackup := StoreRoot
			embedRootBackup := EmbedRoot

			if tt.name == "Test_HandleAsset_embed" {
				StoreRoot = "./not-found"
				EmbedRoot = "embed_test"
				router.Content = fncEMBED
				tt.args.c.URL.Path = "/index.html"
			}

			HandleAsset(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			StoreRoot = storeRootBackup
			EmbedRoot = embedRootBackup

			require.Equal(t, tt.args.want, data, tt.name+" not equal")
		})
	}
}

func Test_WebsocketEcho(t *testing.T) {
	t.Run("WebsocketEcho", func(t *testing.T) {
		serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &router.Context{Request: r, ResponseWriter: w}
			HandleWebsocketEcho(c)
		})
		s := httptest.NewServer(serverHandler)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		w, err := websocket.Dial(u, "", s.URL)
		if err != nil {
			t.Errorf("Dial error %v", err)
		}

		msg := []byte("Hello")

		i, err := w.Write(msg)
		if err != nil {
			t.Errorf("Write error %v", err)
		}

		require.Equal(t, len(msg), i, "WebsocketEcho not equal")
		log.Println(i)
	})
}

func Test_WebsocketChat(t *testing.T) {
	t.Run("WebsocketChat", func(t *testing.T) {
		serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &router.Context{Request: r, ResponseWriter: w}
			HandleWebsocketChat(c)
		})
		s := httptest.NewServer(serverHandler)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		w, err := websocket.Dial(u, "", s.URL)
		if err != nil {
			t.Errorf("Dial error %v", err)
		}

		msg := []byte("Hello")

		i, err := w.Write(msg)
		if err != nil {
			t.Errorf("Write error %v", err)
		}

		require.Equal(t, len(msg), i, "WebsocketChat not equal")
	})
}
