package models

import "github.com/dgrijalva/jwt-go"

type URL struct {
	LongUrl  string `json:"longurl"`
	ShortUrl string `json:"shorturl"`
}

type Claims struct {
	LongUrl string
	jwt.StandardClaims
}

var JwtKey = []byte("my-s3cr3t-k3y")
