package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"router-practice/model"
	"testing"

	"gopkg.in/guregu/null.v4"
)

func TestSetupCookieToken(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		authinfo model.AuthInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "SetupCookieToken",
			args: args{
				w: nil,
				authinfo: model.AuthInfo{
					Name:     null.StringFrom("test_name"),
					IpAddr:   null.StringFrom("192.168.1.1"),
					Platform: null.StringFrom("test_platform"),
					Duration: null.IntFrom(3600),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetupCookieToken(tt.args.w, tt.args.authinfo); (err != nil) != tt.wantErr {
				t.Errorf("SetupCookieToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetClaim(t *testing.T) {
	type args struct {
		w           http.ResponseWriter
		r_cookie    http.Request
		r_header    http.Request
		authinfo    model.AuthInfo
		from_cookie string
		from_header string
	}
	tests := []struct {
		name    string
		args    args
		want    model.AuthInfo
		wantErr bool
	}{
		{
			name: "GetClaim",
			args: args{
				r_cookie: *httptest.NewRequest(http.MethodGet, "/test", nil),
				r_header: *httptest.NewRequest(http.MethodGet, "/test", nil),
				authinfo: model.AuthInfo{
					Name:     null.StringFrom("test_name"),
					IpAddr:   null.StringFrom("192.168.0.1"),
					Platform: null.StringFrom("test_platform"),
					Duration: null.IntFrom(3600),
				},
				from_cookie: "cookie",
				from_header: "header",
			},
			want: model.AuthInfo{
				Name:     null.StringFrom("test_name"),
				IpAddr:   null.StringFrom("192.168.0.1"),
				Platform: null.StringFrom("test_platform"),
				Duration: null.IntFrom(3600),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateKey()
			if err != nil {
				t.Errorf("GenerateKey() error = %v", err)
				return
			}
			token, err := GenerateToken(tt.args.authinfo)
			if err != nil {
				t.Errorf("GenerateToken() error = %v", err)
				return
			}

			tt.args.r_cookie.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})
			tt.args.r_header.Header.Add("Authorization", "Bearer "+token)

			got, err := GetClaim(tt.args.r_cookie, tt.args.from_cookie)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClaim()\nerror = %v\nwants = %v", got, tt.want)
			}
		})
	}
}
