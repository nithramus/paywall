package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../database"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Print("coucou")
		w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("500 - derror"))
	return
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	newUser := database.User{}
	_ = json.Unmarshal(body, &newUser)
	if (newUser.Email == "" || newUser.Password == "") {
		http.Error(w, "Missing Email or Password", http.StatusBadRequest)
	}
	database.UserModel.InsertOne(database.DatabaseCtx, newUser)
	fmt.Fprintf(w, "")
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("500 - derror"))

	

}
