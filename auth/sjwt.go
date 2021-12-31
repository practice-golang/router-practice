package auth

import (
	"log"
	"time"

	"router-practice/model"

	"github.com/brianvoe/sjwt"
)

var Secret string = "secret"

func CreateToken() string {
	// Set Claims
	claims := sjwt.New()
	claims.SetTokenID()
	claims.SetSubject("practice token")
	claims.SetIssuer("practice-golang")
	claims.SetIssuedAt(time.Now())
	claims.SetExpiresAt(time.Now().Add(time.Hour * 24))

	claims.Set("name", "user")

	// Generate jwt
	secretKey := []byte("secret_key_here")
	jwt := claims.Generate(secretKey)

	return jwt
}

func ParseToken(jwt string) {
	claims, err := sjwt.Parse(jwt)
	if err != nil {
		log.Println(err)
	}

	log.Println(claims.Get("name"))

	exp, _ := claims.GetExpiresAt()
	log.Println(exp, time.Now())

	user := model.UserInfo{}
	claims.ToStruct(user)

	log.Println(user)
}