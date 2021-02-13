package main

import (
	"fmt"
	"net/http"

	"./database"
	"./user"
	"github.com/gorilla/mux"
)


func getArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Fprintf(w, "key: "+key)
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/login", user.Login).Methods("POST")
	myRouter.HandleFunc("/signup", user.Signup).Methods("POST")
	myRouter.HandleFunc("/article/{id}", getArticle)
	http.ListenAndServe(":3000", myRouter)
}

func main() {
	client := database.OpenMongoClient()
	defer client.Disconnect(database.DatabaseCtx)
	handleRequest()
}
