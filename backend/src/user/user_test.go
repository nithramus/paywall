package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"../database"
)

// func TestBadLogin(t *testing.T) {

// 	resp, err := http.Post("http://localhost:3000/signup", "application/json", bytes.NewBuffer([]byte("{}")))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode != http.StatusBadRequest {
// 		t.Fatal()
// 	}
// }

// func TestSuccessFullLogin(t *testing.T) {
// 	user := database.User{Email: "baptiste", Password: "test"}
// 	parameters, _ := json.Marshal(user)
// 	resp, err := http.Post("http://localhost:3000/signup", "application/json", bytes.NewBuffer(parameters))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatal()
// 	}
// }

func TestLogin(t *testing.T) {
	user := database.User{Email: "baptiste", Password: "test"}
	parameters, _ := json.Marshal(user)
	http.Post("http://localhost:3000/signup", "application/json", bytes.NewBuffer(parameters))
	resp, err := http.Post("http://localhost:3000/login", "application/json", bytes.NewBuffer(parameters))
	if err != nil {
		t.Fatal()
	}
	ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatal()
	}

}
