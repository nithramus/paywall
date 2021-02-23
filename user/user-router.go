package user

import "github.com/gorilla/mux"

func GetUserRouter(router *mux.Router) *mux.Router {
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/signup", Signup).Methods("POST")
	return router
}
