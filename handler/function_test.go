package handler

import (
	"bytes"
	"embed"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"router-practice/model"
	"router-practice/router"
	"strings"
	"testing"

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

func TestHandleGetDir(t *testing.T) {
	type args struct {
		c        *router.Context
		jsonBody map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test_GetDir_name_asc",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodPost, "/api/dir/list", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				jsonBody: map[string]interface{}{"path": "..", "sort": "name", "order": "asc"},
			},
			want: []byte(`{"path":"..","full-path":"C:\\Users\\high\\Desktop\\pcbangstudio\\workspace\\router-practice","files":[{"name":".git","size":0,"datetime":"2022-01-11 02:07:30","isdir":true},{"name":".github","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"auth","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"embed","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"fd","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"handler","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"html","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"logging","size":0,"datetime":"2022-01-11 02:34:33","isdir":true},{"name":"model","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"router","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"static","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"util","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"variable","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"wsock","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":".gitignore","size":47,"datetime":"2022-01-05 16:34:29","isdir":false},{"name":"Makefile","size":490,"datetime":"2022-01-10 19:44:26","isdir":false},{"name":"README.md","size":803,"datetime":"2022-01-10 19:44:26","isdir":false},{"name":"cover.cmd","size":493,"datetime":"2022-01-05 16:33:22","isdir":false},{"name":"go.mod","size":1210,"datetime":"2022-01-11 02:34:11","isdir":false},{"name":"go.sum","size":8075,"datetime":"2022-01-11 02:34:11","isdir":false},{"name":"main.go","size":756,"datetime":"2022-01-11 00:21:51","isdir":false},{"name":"main_test.go","size":1427,"datetime":"2022-01-10 19:44:26","isdir":false},{"name":"requests.http","size":3445,"datetime":"2022-01-11 03:05:02","isdir":false},{"name":"setup.go","size":3162,"datetime":"2022-01-10 19:44:26","isdir":false}]}`),
		},
		{
			name: "Test_GetDir_name_desc",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodPost, "/api/dir/list", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				jsonBody: map[string]interface{}{"path": "..", "sort": "name", "order": "desc"},
			},
			want: []byte(`{"path":"..","full-path":"C:\\Users\\high\\Desktop\\pcbangstudio\\workspace\\router-practice","files":[{"name":"setup.go","size":3162,"datetime":"2022-01-10 19:44:26","isdir":false},{"name":"requests.http","size":3445,"datetime":"2022-01-11 03:16:22","isdir":false},{"name":"main_test.go","size":1427,"datetime":"2022-01-10 19:44:26","isdir":false},{"name":"main.go","size":756,"datetime":"2022-01-11 00:21:51","isdir":false},{"name":"go.sum","size":8075,"datetime":"2022-01-11 02:34:11","isdir":false},{"name":"go.mod","size":1210,"datetime":"2022-01-11 02:34:11","isdir":false},{"name":"cover.cmd","size":493,"datetime":"2022-01-05 16:33:22","isdir":false},{"name":"README.md","size":803,"datetime":"2022-01-10 19:44:26","isdir":false},{"name":"Makefile","size":490,"datetime":"2022-01-10 19:44:26","isdir":false},{"name":".gitignore","size":47,"datetime":"2022-01-05 16:34:29","isdir":false},{"name":"wsock","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"variable","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"util","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"static","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"router","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"model","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"logging","size":0,"datetime":"2022-01-11 02:34:33","isdir":true},{"name":"html","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"handler","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"fd","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":"embed","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":"auth","size":0,"datetime":"2022-01-10 19:44:26","isdir":true},{"name":".github","size":0,"datetime":"2022-01-04 22:27:22","isdir":true},{"name":".git","size":0,"datetime":"2022-01-11 02:07:30","isdir":true}]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.args.jsonBody)
			tt.args.c.Request = httptest.NewRequest(http.MethodPost, "/api/dir/list", bytes.NewReader(body))
			HandleGetDir(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			var want model.FileList
			err = json.Unmarshal(tt.want, &want)

			var got model.FileList
			err = json.Unmarshal(data, &got)

			require.Equal(t, want.Path, got.Path, "GetDir not equal")
			for i, v := range want.Files {
				require.Equal(t, v.Name, got.Files[i].Name, "GetDir not equal")
			}
		})
	}
}
