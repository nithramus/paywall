package sites

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"paywall/database"
	"paywall/user"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetSites(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value("user").(string)
	var sites []database.Site
	cursor, err := database.SiteModel.Find(database.DatabaseCtx, bson.M{})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(database.DatabaseCtx, &sites); err != nil {
		log.Fatal(err)
	}
	jsonYolo, _ := json.Marshal(sites)
	fmt.Fprintf(w, string(jsonYolo))
}

func AddSite(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	site := database.Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		log.Fatal(err)
	}
	result, err := database.SiteModel.InsertOne(database.DatabaseCtx, site)
	data, _ := json.Marshal(result)
	fmt.Fprintf(w, string(data))
}

func UpdateSite(w http.ResponseWriter, r *http.Request) {

}

func SiteRouter(router *mux.Router) {
	s := router.PathPrefix("/sites").Subrouter()

	s.Use(user.AuthMiddleware)
	s.HandleFunc("", GetSites).Methods("GET")
	s.HandleFunc("", AddSite).Methods("POST")
}
