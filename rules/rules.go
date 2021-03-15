package rules

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"paywall/database"
	"strconv"

	"github.com/gorilla/mux"
)

func GetRules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// accountID := r.Context().Value("accountID").(uint)
	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	var rules []database.AccessRule
	database.Db.Where(&database.AccessRule{SiteID: uint(siteID), Deleted: false}).Find(&rules)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func AddRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	body, _ := ioutil.ReadAll(r.Body)
	var rule database.AccessRule
	err := json.Unmarshal(body, &rule)
	if err != nil {
		log.Fatal(err)
	}
	rule.SiteID = uint(siteID)
	database.Db.Create(&rule)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)

}

func UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	body, _ := ioutil.ReadAll(r.Body)
	var siteMap map[string]interface{}
	err := json.Unmarshal(body, &siteMap)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.ParseUint(vars["ruleID"], 10, 64)
	database.Db.Model(&database.AccessRule{SiteID: uint(siteID), ID: uint(id)}).Updates(&siteMap)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(siteMap)
}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	ruleID, _ := strconv.ParseUint(vars["ruleID"], 10, 64)
	database.Db.Model(&database.AccessRule{SiteID: uint(siteID), ID: uint(ruleID)}).Update("Deleted", true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func AddRuleToOffre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offreID, _ := strconv.ParseUint(vars["offreID"], 10, 64)
	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	ruleID, _ := strconv.ParseUint(vars["ruleID"], 10, 64)
	err := database.Db.Model(&database.Offre{ID: uint(offreID)}).Association("AccessRule").Append([]database.AccessRule{{ID: uint(ruleID), SiteID: uint(siteID)}})
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func RemoveRuleFromOffre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offreID, _ := strconv.ParseUint(vars["offreID"], 10, 64)
	siteID, _ := strconv.ParseUint(vars["siteID"], 10, 64)
	ruleID, _ := strconv.ParseUint(vars["ruleID"], 10, 64)
	err := database.Db.Model(&database.Offre{ID: uint(offreID)}).Association("AccessRule").Delete([]database.AccessRule{{ID: uint(ruleID), SiteID: uint(siteID)}})
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func RuleRouter(router *mux.Router) {
	s := router.PathPrefix("/rules/").Subrouter()

	s.HandleFunc("/site/{siteID}", GetRules).Methods("GET")
	s.HandleFunc("/site/{siteID}", AddRule).Methods("POST")
	// s.HandleFunc("/{ruleID}/site/{siteID}", GetRule).Methods("GET")
	s.HandleFunc("/{ruleID}/site/{siteID}", UpdateRule).Methods("PUT")
	s.HandleFunc("/{ruleID}/site/{siteID}", DeleteRule).Methods("DELETE")
	s.HandleFunc("/{ruleID}/site/{siteID}/offre/{offreID}", AddRuleToOffre).Methods("POST")
	s.HandleFunc("/{ruleID}/site/{siteID}/offre/{offreID}", RemoveRuleFromOffre).Methods("DELETE")
}
