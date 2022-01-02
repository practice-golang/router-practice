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
	Name     null.String `json:"name"`
	IpAddr   null.String `json:"ip-addr"`
	Platform null.String `json:"platform"`
	Duration null.Int    `json:"duration"`
}
