package handlers

import (
	"encoding/json"
	"fmt"
	"go-module/helpers"
	"go-module/models"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func A(w http.ResponseWriter, r *http.Request) {
	// get data from frontend team
	var res RESP
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		log.Println("Err decode is : ", err)
	}

	// check domain name, if ok then create short url
	Re := regexp.MustCompile("^(https://|http://)?[a-zA-Z0-9\\./:]+(/|\\.)[a-zA-Z0-9]+$")
	if !Re.MatchString(res.Value) {
		msg := models.URL{LongUrl: res.Value, ShortUrl: "toang"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&msg)

		fmt.Println("Failed")
		return
	}
	shortU := uuid.NewV4().String()
	bshortU := []byte(shortU)[:8]
	shortU = string(bshortU)

	// fix domain name without "http" or "https"
	if !helpers.Exist(res.Value, "http://") && !helpers.Exist(res.Value, "https://") {
		res.Value = "http://" + res.Value
	}

	// connect to the database
	db, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Println("Err db is : ", err)
	}
	defer db.Close()

	// create short url for this long url, also jwt
	u1 := models.URL{LongUrl: res.Value, ShortUrl: shortU}
	db.DB("orlab").C("urls").Insert(&u1)
	fmt.Println("Long Url is : ", res.Value)
	fmt.Println("Short Url is : ", shortU)

	// give cookie
	helpers.CreateJWT(w, r, u1.LongUrl, 30*time.Minute)

	// give response back to frontend side
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&u1)
}

func GiveLink(w http.ResponseWriter, r *http.Request) {
	// verify JWT
	u := helpers.VerifyJWT(w, r)
	if u.LongUrl == "" {
		return
	}

	// get data from url
	vars := mux.Vars(r)
	name := vars["name"]

	// open database
	db, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Println("Err db is : ", err)
	}
	defer db.Close()

	// check name in database, then get the long url
	var result models.URL
	err1 := db.DB("orlab").C("urls").Find(bson.M{"shorturl": name}).One(&result)
	if err1 != nil {
		log.Println("Wrong link!")
		return
	}
	http.Redirect(w, r, result.LongUrl, http.StatusSeeOther)
}

// helpers

type RESP struct {
	Value string
}
