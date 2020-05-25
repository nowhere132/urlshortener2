package helpers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-module/models"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Exist(x, y string) bool {
	if len(x) < len(y) {
		return false
	}

	for i := 0; i < len(y); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func CreateJWT(w http.ResponseWriter, r *http.Request, RealUrl string, t time.Duration) {
	expirationTime := time.Now().Add(t)

	claim := models.Claims{
		LongUrl: RealUrl,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(models.JwtKey)
	if err != nil {
		log.Println(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	fmt.Println("ok jwt created successfully")
}

func VerifyJWT(w http.ResponseWriter, r *http.Request) models.URL {
	c, err := r.Cookie("token")
	if err != nil {
		log.Println("Err cookie is : ", err)
		return models.URL{}
	}

	tokenString := c.Value
	claim := models.Claims{}

	token, err1 := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})
	if err1 != nil {
		log.Println("Err parse is : ", err1)
		return models.URL{}
	}
	if !token.Valid {
		log.Println(http.StatusUnauthorized)
		return models.URL{}
	}

	// connect to the database
	db, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Println("Err db is : ", err)
	}
	defer db.Close()

	var answer models.URL
	db.DB("orlab").C("urls").Find(bson.M{"longurl": claim.LongUrl}).One(&answer)
	return answer
}
