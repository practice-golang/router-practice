package auth

import (
	"errors"
	"log"
	"net/http"
	"router-practice/model"
	"strings"
)

func SetupCookieToken(w http.ResponseWriter, authinfo model.AuthInfo) error {
	token, err := GenerateToken(authinfo)
	if err != nil {
		log.Println(err)
		return err
	}

	SetCookieHeader(w, token, authinfo.Duration.Int64)

	return nil
}

func GetClaim(r http.Request, from string) (model.AuthInfo, error) {
	var result model.AuthInfo
	var dataCookie *http.Cookie
	var dataHeader string
	var token string
	var err error

	switch from {
	case "cookie":
		dataCookie, err = r.Cookie("token")
		if err != nil {
			log.Println("GetCookie cookie:", err)
			return result, err
		}

		token = dataCookie.Value
	case "header":
		dataHeader = r.Header.Get("Authorization")
		dataHeaders := strings.Split(dataHeader, " ") // Bearer token
		if dataHeaders[0] != "Bearer" {
			log.Println("GetCookie cookie:", "Bearer not found")
			return result, errors.New("bearer not found")
		}

		token = dataHeaders[1]
	}

	_, result, err = ParseToken(token)
	if err != nil {
		// log.Println("GetCookie parse token:", err)
		return model.AuthInfo{}, err
	}

	return result, nil
}
