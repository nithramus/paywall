package offres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"paywall/database"
	"paywall/test"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetOffres(t *testing.T) {
	client := test.SignupLogin()
	resp, err := client.Get("http://localhost:3001/offres")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code: ", resp.StatusCode)
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	offre := database.Offre{Name: "propro"}
	params, _ := json.Marshal(offre)
	resp, err = client.Post("http://localhost:3001/offres", "application/json", bytes.NewBuffer(params))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code: ", resp.StatusCode)
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]string

	json.Unmarshal(body, &result)
	fmt.Println(result["InsertedID"])
	id, _ := primitive.ObjectIDFromHex(result["InsertedID"])
	offre = database.Offre{ID: id, Name: "##############"}
	params, _ = json.Marshal(offre)
	req, err := http.NewRequest("PUT", "http://localhost:3001/offres/"+result["InsertedID"], bytes.NewReader(params))

	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(body)

	req, err = http.NewRequest("DELETE", "http://localhost:3001/offres/"+result["InsertedID"], nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)

}
