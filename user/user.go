package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"paywall/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Println(vars)
		type MyCustomClaims struct {
			Id string `json:"id"`
			jwt.StandardClaims
		}
		var bearToken string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "token" {
				bearToken = cookie.Value
			}
		}
		claims := &MyCustomClaims{}
		tkn, err := jwt.ParseWithClaims(bearToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("eirueiztuiretuire"), nil
		})
		if err != nil {
			fmt.Println(err)
			if err == jwt.ErrSignatureInvalid {
				fmt.Print("etresz")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			fmt.Println("unauth")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.Id)

		r = r.WithContext(ctx)
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
	expTime := time.Now().Add(time.Minute * 60 * 24 * 300)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = user.ID
	atClaims["email"] = user.Email
	atClaims["exp"] = expTime.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("eirueiztuiretuire"))
	if err != nil {
		panic(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	})

}