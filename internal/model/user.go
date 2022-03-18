package model

import (
	"gopkg.in/guregu/null.v4"
)

type UserInfo struct {
	Name null.String `json:"name"`
	Age  null.Int    `json:"age"`
}

type SignIn struct {
	Name     null.String `json:"name"`
	Password null.String `json:"password"`
}

type AuthInfo struct {
	Name     null.String `json:"name"     mapstructure:"name"`
	IpAddr   null.String `json:"ip-addr"  mapstructure:"ip-addr"`
	Platform null.String `json:"platform" mapstructure:"platform"`
	Duration null.Int    `json:"duration" mapstructure:"duration"`
}
