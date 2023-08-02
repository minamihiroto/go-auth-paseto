package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/o1egl/paseto"
)

var symmetricKey = []byte("YELLOW SUBMARINE, BLACK WIZARDRY") // Must be 32 bytes

func issueHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)

	jsonToken := paseto.JSONToken{
		Audience:   "test",
		Issuer:     "test_service",
		Jti:        "123",
		Subject:    "test_subject",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  now,
	}
	jsonToken.Set("data", "this is a signed message")
	footer := "some footer"

	v2 := paseto.NewV2()

	token, err := v2.Encrypt(symmetricKey, jsonToken, footer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, token)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")

	if len(splitToken) != 2 {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	reqToken := splitToken[1]

	var newJsonToken paseto.JSONToken
	var newFooter string

	v2 := paseto.NewV2()

	err := v2.Decrypt(reqToken, symmetricKey, &newJsonToken, &newFooter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, "Success!!!!!")
}

func main() {
	http.HandleFunc("/issue", issueHandler)
	http.HandleFunc("/verify", verifyHandler)

	fmt.Println("Starting server at port 8080")
	http.ListenAndServe(":8080", nil)
}

//go run main.go
//curl -X GET http://127.0.0.1:8080/issue
//curl -X GET -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/verify