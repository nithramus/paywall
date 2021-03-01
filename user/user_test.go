package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"paywall/database"
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
	resp, err := http.Post("http://localhost:3000/signup", "application/json", bytes.NewBuffer(parameters))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Wrong status code")
	}
	resp, err = http.Post("http://localhost:3000/login", "application/json", bytes.NewBuffer(parameters))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Wrong status code")
	}
	var jwtToken string
	for _, cookie := range resp.Cookies() {
		fmt.Println("Found a cookie named:", cookie.Name, cookie.Value)
		jwtToken = cookie.Value
	}
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookie := &http.Cookie{
		Name:  "token",
		Value: jwtToken,
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse("http://localhost:300")
	jar.SetCookies(u, cookies)
	client := &http.Client{
		Jar: jar,
	}
	client.Get("http://localhost:3000/articles/list")
	ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Wrong status code")
	}

}
