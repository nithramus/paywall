package user

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {

	resp, err := http.Post("http://localhost:3000/signup", "application/json", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		log.Fatal(err)
	}
	ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatal()
	}
}