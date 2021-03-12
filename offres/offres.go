package offres

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

func GetOffres(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value("accountID").(uint)
	var offres []database.Offre
	offres = make([]database.Offre, 0)

	database.Db.Where(&database.Offre{Deleted: false, AccountID: accountID}).Preload("Sites").Find(&offres)

	jsonYolo, err := json.Marshal(offres)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(jsonYolo))
}

func AddOffre(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	offre := database.Offre{}
	err = json.Unmarshal(body, &offre)
	offre.AccountID = r.Context().Value("accountID").(uint)
	if err != nil {
		log.Fatal(err)
	}
	result := database.Db.Create(&offre)
	if result.Error != nil {
		log.Fatal(err)
	}
	data, _ := json.Marshal(offre)
	fmt.Fprintf(w, string(data))
}

func UpdateOffre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	offre := database.Offre{}
	err = json.Unmarshal(body, &offre)
	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.ParseUint(vars["offreID"], 10, 64)
	accountID := r.Context().Value("accountID").(uint)
	database.Db.Model(&database.Offre{AccountID: accountID, ID: uint(id)}).Updates(&offre)

	if err != nil {
		log.Fatal(err)
	}
	data, _ := json.Marshal(&offre)
	w.Write(data)
}

func DeleteOffre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["offreID"], 10, 64)
	accountID := r.Context().Value("accountID").(uint)
	result := database.Db.Model(&database.Offre{AccountID: accountID, ID: uint(id)}).Update("Deleted", true)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	data, _ := json.Marshal(result)
	w.Write(data)
}

func OffreRouter(router *mux.Router) {
	s := router.PathPrefix("/offres").Subrouter()

	s.HandleFunc("", GetOffres).Methods("GET")
	s.HandleFunc("", AddOffre).Methods("POST")
	s.HandleFunc("/{offreID}", UpdateOffre).Methods("PUT")
	s.HandleFunc("/{offreID}", DeleteOffre).Methods("DELETE")
}
