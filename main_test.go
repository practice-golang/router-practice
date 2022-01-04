package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"router-practice/logging"
	"router-practice/router"
	"router-practice/variable"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
	type args struct{ c *router.Context }
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test_main",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("Hello world GET"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Address = "localhost"
			// Port = "4416"
			// go main()

			// res, err := http.Get("http://" + Address + ":" + Port + "/hello")
			// if err != nil {
			// 	t.Fatal("http.Get", err)
			// }
			// defer res.Body.Close()

			doSetup()

			ServerHandler.ServeHTTP(tt.args.c.ResponseWriter, tt.args.c.Request)
			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal("read body", err)
			}
			fname := variable.ProgramName + "-" + time.Now().Format("20060102") + ".log"
			logging.F.Close()
			os.Remove(fname)

			require.Equal(t, tt.want, data, "not equal")
		})
	}
}
