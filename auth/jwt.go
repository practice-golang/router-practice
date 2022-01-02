package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"router-practice/model"
	"router-practice/variable"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mitchellh/mapstructure"
)

var Secret string = "secret"
var Alg jwa.SignatureAlgorithm = jwa.RS384
var KeySET jwk.Set
var RealKey jwk.Key

// GenerateKey - 키 생성
func GenerateKey() {
	// log.Println("alg/key:", Alg, Secret)

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("failed to generate private key: %s\n", err)
		return
	}

	pubKey, err := jwk.New(privKey.PublicKey)
	if err != nil {
		fmt.Printf("failed to create JWK: %s\n", err)
		return
	}

	pubKey.Set(jwk.AlgorithmKey, Alg)
	pubKey.Set(jwk.KeyIDKey, Secret)

	bogusKey := jwk.NewSymmetricKey()
	bogusKey.Set(jwk.AlgorithmKey, jwa.NoSignature)
	bogusKey.Set(jwk.KeyIDKey, "otherkey")

	KeySET = jwk.NewSet()
	KeySET.Add(pubKey)
	KeySET.Add(bogusKey)

	RealKey, err = jwk.New(privKey)
	if err != nil {
		log.Printf("failed to create JWK: %s\n", err)
		return
	}

	RealKey.Set(jwk.KeyIDKey, Secret)
	RealKey.Set(jwk.AlgorithmKey, Alg)
}

// GenerateToken - 토큰 생성
func GenerateToken(authinfo model.AuthInfo) (string, error) {
	token, err := jwt.NewBuilder().
		Issuer(variable.ProgramName).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(time.Duration(authinfo.Duration.Int64)*time.Second)).
		Subject("auth token").
		Claim("token", authinfo).
		Build()

	if err != nil {
		log.Printf("failed to begin to build: %s\n", err)
		return "", err
	}

	// token.Set("token", authinfo)

	signed, err := jwt.Sign(token, Alg, RealKey)
	if err != nil {
		log.Printf("failed to generate signed payload: %s\n", err)
		return "", err
	}

	result := string(signed)

	return result, err
}

// ParseToken - 토큰 파싱
func ParseToken(payloadSTR string) (jwt.Token, model.AuthInfo, error) {
	payload := []byte(payloadSTR)

	token, err := jwt.Parse(
		payload,
		jwt.WithKeySet(KeySET),
	)
	if err != nil {
		log.Printf("parse payload: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	var authinfo model.AuthInfo

	cfg := &mapstructure.DecoderConfig{
		Result:     &authinfo,
		DecodeHook: ConvertToNullTypeHookFunc,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		log.Printf("decoder set: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	claim, valid := token.Get("token")
	if !valid {
		log.Printf("token to claim: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	err = decoder.Decode(claim)
	if err != nil {
		log.Printf("decode claim to struct: %s\n", err)
		return nil, model.AuthInfo{}, err
	}

	now := time.Now()
	if token.Expiration().Before(now) {
		log.Printf("token expired: %s\n", err)
		return nil, model.AuthInfo{}, errors.New("token expired")
	}

	return token, authinfo, err
}
