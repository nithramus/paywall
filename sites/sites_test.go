package sites

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"paywall/database"
	"paywall/test"
	"strconv"
	"testing"
)

func TestSite(t *testing.T) {
	client := test.SignupLogin()
	newSite := database.Site{Name: "Website Name"}
	params, _ := json.Marshal(newSite)
	resp, err := client.Post("http://localhost:3001/sites", "application/json", bytes.NewBuffer(params))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code: ", resp.StatusCode)
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &newSite)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.Get("http://localhost:3001/sites")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code: ", resp.StatusCode)
		t.Fatal(err)
	}
	var sites []database.Site
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &sites)
	newSite.Name = "tgestst"
	params, _ = json.Marshal(newSite)
	req, err := http.NewRequest("PUT", "http://localhost:3001/sites/"+strconv.FormatUint(uint64(newSite.ID), 10), bytes.NewReader(params))
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code: ", resp.StatusCode)
		t.Fatal(err)
	}
	req, err = http.NewRequest("DELETE", "http://localhost:3001/sites/"+strconv.FormatUint(uint64(newSite.ID), 10), nil)
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code: ", resp.StatusCode)
		t.Fatal(err)
	}

}
