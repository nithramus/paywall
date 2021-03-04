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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSites(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(string)
	var sites []database.Site
	cursor, err := database.SiteModel.Find(database.DatabaseCtx, bson.M{"deleted": false, "userId": userID})
	if err != nil {
		log.Fatal(err)
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
	site.UserID = r.Context().Value("userId").(string)
	if err != nil {
		log.Fatal(err)
	}
	result, err := database.SiteModel.InsertOne(database.DatabaseCtx, site)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*result)
	data, _ := json.Marshal(*result)
	fmt.Fprintf(w, string(data))
}

func UpdateSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	site := database.Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		log.Fatal(err)
	}
	update := bson.M{"$set": site}

	id, _ := primitive.ObjectIDFromHex(vars["siteId"])
	userId := r.Context().Value("userId").(string)

	result, err := database.SiteModel.UpdateOne(database.DatabaseCtx, bson.M{"_id": id, "userId": userId}, update)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := json.Marshal(result)
	w.Write(data)
}

func DeleteSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["siteId"])
	site := database.Site{Deleted: true}
	result, err := database.SiteModel.UpdateOne(database.DatabaseCtx, bson.M{"_id": id}, bson.M{"$set": site})
	if err != nil {
		log.Fatal(err)
	}
	data, _ := json.Marshal(result)
	w.Write(data)
}

func SiteRouter(router *mux.Router) {
	s := router.PathPrefix("/sites").Subrouter()

	s.Use(user.AuthMiddleware)
	s.HandleFunc("", GetSites).Methods("GET")
	s.HandleFunc("", AddSite).Methods("POST")
	s.HandleFunc("/{siteId}", UpdateSite).Methods("PUT")
	s.HandleFunc("/{siteId}", DeleteSite).Methods("DELETE")
}
