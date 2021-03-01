package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"paywall/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearToken := r.Header.Get("token")
		fmt.Println(bearToken)
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(bearToken, claims, func(token *jwt.Token) (interface{}, error) {
			return "eirueiztuiretuire", nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// http.Error(w, "Invalid Auth", http.StatusUnauthorized)

		next.ServeHTTP(w, r)
	})
}

func Signup(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	newUser := database.User{}
	_ = json.Unmarshal(body, &newUser)
	if newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "Missing Email or Password", http.StatusBadRequest)
		return
	}
	pass := []byte(newUser.Password)
	hashedPass, errr := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if errr != nil {
		log.Fatal(errr)
		http.Error(w, "Fail to hash password", http.StatusBadRequest)
	}
	newUser.Password = string(hashedPass)
	_, err = database.UserModel.InsertOne(database.DatabaseCtx, newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user database.User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	loggingUser := database.User{}
	err = json.Unmarshal(body, &loggingUser)
	if err != nil {
		log.Fatal(err)
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(loggingUser.Password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	err = database.UserModel.FindOne(database.DatabaseCtx, bson.M{"email": loggingUser.Email}).Decode(&user)

	if err != nil {
		panic(err)
	}

	if string(hashedPass) != user.Password && loggingUser.Email != user.Email {
		http.Error(w, "Bad password or email", http.StatusBadRequest)
		return
	}
	// Create the token
	expTime := time.Now().Add(time.Minute * 15)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = user.Email
	atClaims["exp"] = expTime.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("eirueiztuiretuire"))
	fmt.Println(token)
	if err != nil {
		panic(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	})

}
