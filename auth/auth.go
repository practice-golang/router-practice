package auth

import (
	"log"
	"net/http"
	"router-practice/model"
)

func SetupToken(w http.ResponseWriter, authinfo model.AuthInfo) error {
	token, err := GenerateToken(authinfo)
	if err != nil {
		log.Println(err)
		return err
	}

	SetCookieHeader(w, token, authinfo.Duration.Int64)

	return nil
}

func GetClaim(r http.Request) (model.AuthInfo, error) {
	var result model.AuthInfo

	token, err := r.Cookie("token")
	if err != nil {
		log.Println("GetCookie cookie:", err)
		return result, err
	}

	_, result, err = ParseToken(token.Value)
	if err != nil {
		log.Println("GetCookie parse token:", err)
		return model.AuthInfo{}, err
	}

	return result, nil
}
