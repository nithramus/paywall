package offres

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

func GetOffres(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(string)
	var offres []database.Offre
	offres = make([]database.Offre, 0)
	cursor, err := database.OffreModel.Find(database.DatabaseCtx, bson.M{"deleted": false, "userId": userID})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(database.DatabaseCtx, &offres); err != nil {
		log.Fatal(err)
	}
	jsonYolo, err := json.Marshal(offres)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonYolo), offres)

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
	offre.UserID = r.Context().Value("userId").(string)
	if err != nil {
		log.Fatal(err)
	}
	result, err := database.OffreModel.InsertOne(database.DatabaseCtx, offre)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*result)
	data, _ := json.Marshal(*result)
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
	update := bson.M{"$set": offre}

	id, _ := primitive.ObjectIDFromHex(vars["offreId"])
	userId := r.Context().Value("userId").(string)

	result, err := database.OffreModel.UpdateOne(database.DatabaseCtx, bson.M{"_id": id, "userId": userId}, update)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := json.Marshal(result)
	w.Write(data)
}

func DeleteOffre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["offreId"])
	offre := database.Offre{Deleted: true}
	result, err := database.OffreModel.UpdateOne(database.DatabaseCtx, bson.M{"_id": id}, bson.M{"$set": offre})
	if err != nil {
		log.Fatal(err)
	}
	data, _ := json.Marshal(result)
	w.Write(data)
}

func OffreRouter(router *mux.Router) {
	s := router.PathPrefix("/offres").Subrouter()

	s.Use(user.AuthMiddleware)
	s.HandleFunc("", GetOffres).Methods("GET")
	s.HandleFunc("", AddOffre).Methods("POST")
	s.HandleFunc("/{offreId}", UpdateOffre).Methods("PUT")
	s.HandleFunc("/{offreId}", DeleteOffre).Methods("DELETE")
}
