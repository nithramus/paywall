package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"./database"
	"./user"
	"github.com/gorilla/mux"
)

func getArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Printf(key, "test")
}

func RecoverWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		defer func() {
			r := recover()
			if r != nil {

				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}

				_, fn, line, _ := runtime.Caller(2)
				log.Printf("[error] %s:%d %v", fn, line, err)

				http.Error(w, "Internal serveur error", http.StatusInternalServerError)
			}

		}()
		next.ServeHTTP(w, r)
	})
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(RecoverWrap)
	user.GetUserRouter(myRouter)

	myRouter.Use(user.AuthMiddleware)
	myRouter.HandleFunc("/articles", getArticle)
	mw := LogMiddleware(myRouter)
	http.ListenAndServe(":3000", mw)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client := database.OpenMongoClient()
	defer client.Disconnect(database.DatabaseCtx)
	handleRequest()
}
