package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"paywall/database"
	"paywall/sites"
	"paywall/user"

	"github.com/gorilla/mux"
)

func getArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Printf(key, "test")
	fmt.Fprintf(w, "nike la police")
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
				_, fn, line, _ = runtime.Caller(0)
				log.Printf("[error] %s:%d %v", fn, line, err)
				http.Error(w, "Internal serveur error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LogMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		fmt.Println(vars)
		fmt.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func handleRequest() {
	myRouter := mux.NewRouter()
	myRouter.Use(RecoverWrap)
	user.GetUserRouter(myRouter)
	sites.SiteRouter(myRouter)
	mw := LogMiddleware(myRouter)
	err := http.ListenAndServe(":3001", mw)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client := database.OpenMongoClient()
	defer client.Disconnect(database.DatabaseCtx)
	handleRequest()
}
