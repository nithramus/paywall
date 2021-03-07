package sites

import (
	"net/http"

	"github.com/gorilla/mux"
)

func AddSite(w http.ResponseWriter, r *http.Request) {

}

func SiteRouter(router *mux.Router) {
	s := router.PathPrefix("/sites").Subrouter()
	s.HandleFunc("", AddSite).Methods("GET")
}
