package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"paywall/database"
)

func SignupLogin() *http.Client {
	user := database.User{Email: "baptiste@test.fr", Password: "testtest"}
	parameters, _ := json.Marshal(user)
	resp, _ := http.Post("http://localhost:3001/signup", "application/json", bytes.NewBuffer(parameters))
	resp, _ = http.Post("http://localhost:3001/login", "application/json", bytes.NewBuffer(parameters))
	var jwtToken string
	for _, cookie := range resp.Cookies() {
		jwtToken = cookie.Value
	}
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookie := &http.Cookie{
		Name:  "token",
		Value: jwtToken,
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse("http://localhost:3001")
	jar.SetCookies(u, cookies)
	client := &http.Client{
		Jar: jar,
	}
	return client
}
