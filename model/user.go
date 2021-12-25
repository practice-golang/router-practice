package model

import (
	"gopkg.in/guregu/null.v4"
)

type UserInfo struct {
	Name null.String `json:"name"`
	Age  null.Int    `json:"age"`
}
