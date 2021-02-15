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

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Prints logs")
		fmt.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	user.GetUserRouter(myRouter)
	// s := myRouter.PathPrefix("/article").Subrouter()
	// s.Use(user.AuthMiddleware)
	// s.HandleFunc("/{id}", getArticle)
	mw := LogMiddleware(myRouter)
	http.ListenAndServe(":3000", mw)
}

func main() {
	client := database.OpenMongoClient()
	defer client.Disconnect(database.DatabaseCtx)
	handleRequest()
}
