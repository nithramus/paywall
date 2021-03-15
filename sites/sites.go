package sites

import (
	"encoding/json"
	"fmt"
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
	database.Db.Where(&database.Site{AccountID: accountID}).Preload("Offres").Find(&sites)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func GetSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID := r.Context().Value("accountID").(uint)
	id, _ := strconv.ParseUint(vars["siteID"], 10, 64)

	var site database.Site
	database.Db.Where(&database.Site{AccountID: accountID, ID: uint(id)}).Preload("Offres").First(&site)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
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
	offre := database.Offre{Name: "Default", AccountID: accountID}
	database.Db.Create(&offre)
	accessRule := database.AccessRule{SiteID: site.ID, Name: "Default"}
	database.Db.Create(&accessRule)

	err = database.Db.Model(&site).Association("Offres").Append(&offre)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Db.Model(&offre).Association("AccessRules").Append(&accessRule)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func UpdateSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID := r.Context().Value("accountID").(uint)
	// var site database.Site
	body, _ := ioutil.ReadAll(r.Body)
	var siteMap map[string]interface{}
	err := json.Unmarshal(body, &siteMap)
	fmt.Println(siteMap)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	database.Db.Model(&database.Site{AccountID: accountID, ID: uint(id)}).Updates(&siteMap)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(siteMap)
}

func DeleteSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID := r.Context().Value("accountID").(uint)
	id, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	database.Db.Model(&database.Site{AccountID: accountID, ID: uint(id)}).Update("Deleted", true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func AddOffreToSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offreID, _ := strconv.ParseUint(vars["offreID"], 10, 64)
	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	fmt.Println(siteID, offreID)
	err := database.Db.Model(&database.Site{ID: uint(siteID)}).Association("Offres").Append([]database.Offre{{ID: uint(offreID)}})
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func RemoveOffreFromSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offreID, _ := strconv.ParseUint(vars["offreID"], 10, 64)
	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	fmt.Println(siteID, offreID)
	err := database.Db.Model(&database.Site{ID: uint(siteID)}).Association("Offres").Delete([]database.Offre{{ID: uint(offreID)}})
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func SiteRouter(router *mux.Router) {
	s := router.PathPrefix("/sites").Subrouter()

	s.HandleFunc("", GetSites).Methods("GET")
	s.HandleFunc("", AddSite).Methods("POST")
	s.HandleFunc("/{siteID}", GetSite).Methods("GET")
	s.HandleFunc("/{siteID}", UpdateSite).Methods("PUT")
	s.HandleFunc("/{siteID}", DeleteSite).Methods("DELETE")
	s.HandleFunc("/{siteID}/offre/{offreID}", AddOffreToSite).Methods("POST")
	s.HandleFunc("/{siteID}/offre/{offreID}", RemoveOffreFromSite).Methods("DELETE")
}
