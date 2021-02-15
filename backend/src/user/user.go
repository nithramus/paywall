package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../database"
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
		http.Error(w, "Invalid Auth", http.StatusUnauthorized)
		return
		// next.ServeHTTP(w, r)
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
	fmt.Printf(string(hashedPass))
	newUser.Password = string(hashedPass)
	_, err = database.UserModel.InsertOne(database.DatabaseCtx, newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("user added")
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
	err = database.UserModel.FindOne(database.DatabaseCtx, bson.M{"email": loggingUser.Email}).Decode(user)
	if err != nil {
		log.Fatal(err)
	}
	if string(hashedPass) != user.Password && loggingUser.Email != user.Email {
		log.Fatal(err)
		http.Error(w, "bad password or email", http.StatusBadRequest)
		return
	}
	// expirationTime := time.Now().Add(5 * time.Minute)
	// claims := &Claims{
	// 	Username: user.Email,
	// 	StandardClaims: jwt.StandardClaims{
	// 		// In JWT, the expiry time is expressed as unix milliseconds
	// 		ExpiresAt: expirationTime.Unix(),
	// 	},
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenString, err := token.SignedString("etfetreztr")
	// if err != nil {
	// 	// If there is an error in creating the JWT return an internal server error
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })
}
