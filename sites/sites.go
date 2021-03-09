package sites

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"paywall/database"
	"strconv"

	"github.com/gorilla/mux"
)

func GetSites(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value("accountID").(uint)
	var sites []database.Site
	database.Db.Where(&database.Site{AccountID: accountID}).Find(&sites)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func AddSite(w http.ResponseWriter, r *http.Request) {
	var site database.Site

	accountID := r.Context().Value("accountID").(uint)
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &site)
	if err != nil {
		log.Fatal(err)
	}
	site.AccountID = accountID
	database.Db.Create(&site)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func UpdateSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID := r.Context().Value("accountID").(uint)
	var site database.Site
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &site)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	database.Db.Model(&database.Site{AccountID: accountID, ID: uint(id)}).Updates(&site)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func DeleteSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID := r.Context().Value("accountID").(uint)
	id, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	database.Db.Model(&database.Site{AccountID: accountID, ID: uint(id)}).Update("Deleted", true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func SiteRouter(router *mux.Router) {
	s := router.PathPrefix("/sites").Subrouter()

	s.HandleFunc("", GetSites).Methods("GET")
	s.HandleFunc("", AddSite).Methods("POST")
	s.HandleFunc("/{siteID}", UpdateSite).Methods("PUT")
	s.HandleFunc("/{siteID}", DeleteSite).Methods("DELETE")
}
