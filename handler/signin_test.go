package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"router-practice/router"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Signin(t *testing.T) {
	type args struct {
		c        *router.Context
		want     []byte
		run_func router.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_Signin",
			args: args{
				c: &router.Context{
					Request: httptest.NewRequest(
						"GET", "/signin",
						bytes.NewBuffer([]byte(`{"name": "test_user","password": "12345"}`))),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want:     []byte(`"Signin success"`),
				run_func: Signin,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.run_func(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "Not equal")
		})
	}
}

func Test_SigninAPI(t *testing.T) {
	type args struct {
		c        *router.Context
		want     []byte
		run_func router.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_SigninAPI",
			args: args{
				c: &router.Context{
					Request: httptest.NewRequest(
						"GET", "/signin",
						bytes.NewBuffer([]byte(`{"name": "test_user","password": "12345"}`))),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want:     []byte("Signin success"),
				run_func: SigninAPI,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.run_func(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			result := make(map[string]string)
			err = json.Unmarshal(data, &result)
			if err != nil {
				t.Errorf("expected error to parse data %v", err)
			}

			require.Equal(t, tt.args.want, []byte(result["msg"]), "Not equal")
		})
	}
}
