package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

var sampleSecretKey = []byte("secret")

func main() {
	http.HandleFunc("/home", handlePage)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["authorized"] = true
	claims["user"] = "John Doe"

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Message: ", message)

	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
