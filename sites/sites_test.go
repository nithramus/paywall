package sites

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"paywall/test"
	"testing"
)

func TestGetSites(t *testing.T) {
	client := test.SignupLogin()
	resp, err := client.Get("http://localhost:3001/sites")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code: ", resp.StatusCode)
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

}
