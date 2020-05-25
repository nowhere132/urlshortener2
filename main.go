package main

import (
	"go-module/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter()

	// apis
	myRouter.HandleFunc("/a", handlers.A).Methods("POST")

	// handlers
	myRouter.HandleFunc("/s/{name}", handlers.GiveLink)

	myRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates/")))
	http.ListenAndServe(":8000", myRouter)
}
